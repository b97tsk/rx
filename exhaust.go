package rx

import (
	"sync"
	"sync/atomic"
)

// ExhaustAll flattens a higher-order Observable into a first-order Observable
// by dropping inner Observables while the previous inner Observable has not
// yet completed.
func ExhaustAll[_ Observable[T], T any]() Operator[Observable[T], T] {
	return ExhaustMap(identity[Observable[T]])
}

// ExhaustMapTo converts the source Observable into a higher-order Observable,
// by mapping each source value to the same Observable, then flattens it into
// a first-order Observable using ExhaustAll.
func ExhaustMapTo[T, R any](inner Observable[R]) Operator[T, R] {
	return ExhaustMap(func(T) Observable[R] { return inner })
}

// ExhaustMap converts the source Observable into a higher-order Observable,
// by mapping each source value to an Observable, then flattens it into
// a first-order Observable using ExhaustAll.
func ExhaustMap[T, R any](mapping func(v T) Observable[R]) Operator[T, R] {
	return NewOperator(
		func(source Observable[T]) Observable[R] {
			return exhaustMapObservable[T, R]{source, mapping}.Subscribe
		},
	)
}

type exhaustMapObservable[T, R any] struct {
	Source  Observable[T]
	Mapping func(T) Observable[R]
}

func (obs exhaustMapObservable[T, R]) Subscribe(c Context, sink Observer[R]) {
	c, cancel := c.WithCancel()
	sink = sink.DoOnTermination(cancel)

	var x struct {
		Context  atomic.Value
		Complete atomic.Bool
		Worker   struct {
			sync.WaitGroup
		}
	}

	x.Context.Store(c.Context)

	startWorker := func(v T) {
		obs1 := obs.Mapping(v)
		w, cancelw := c.WithCancel()

		x.Context.Store(w.Context)
		x.Worker.Add(1)

		obs1.Subscribe(w, func(n Notification[R]) {
			switch n.Kind {
			case KindNext:
				sink(n)

			case KindError, KindComplete:
				defer x.Worker.Done()

				cancelw()

				switch n.Kind {
				case KindError:
					if x.Context.Swap(sentinel) != sentinel {
						sink(n)
					}
				case KindComplete:
					if x.Context.CompareAndSwap(w.Context, c.Context) && x.Complete.Load() && x.Context.CompareAndSwap(c.Context, sentinel) {
						sink(n)
					}
				}
			}
		})
	}

	obs.Source.Subscribe(c, func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
			if x.Context.Load() == c.Context {
				startWorker(n.Value)
			}

		case KindError:
			old := x.Context.Swap(sentinel)

			cancel()
			x.Worker.Wait()

			if old != sentinel {
				sink.Error(n.Error)
			}

		case KindComplete:
			x.Complete.Store(true)

			if x.Context.CompareAndSwap(c.Context, sentinel) {
				sink.Complete()
			}
		}
	})
}
