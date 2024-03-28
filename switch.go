package rx

import (
	"sync"
	"sync/atomic"
)

// SwitchAll flattens a higher-order Observable into a first-order Observable
// by subscribing to only the most recently emitted of those inner Observables.
func SwitchAll[_ Observable[T], T any]() Operator[Observable[T], T] {
	return SwitchMap(identity[Observable[T]])
}

// SwitchMapTo converts the source Observable into a higher-order Observable,
// by mapping each source value to the same Observable, then flattens it into
// a first-order Observable using SwitchAll.
func SwitchMapTo[T, R any](inner Observable[R]) Operator[T, R] {
	return SwitchMap(func(T) Observable[R] { return inner })
}

// SwitchMap converts the source Observable into a higher-order Observable,
// by mapping each source value to an Observable, then flattens it into
// a first-order Observable using SwitchAll.
func SwitchMap[T, R any](mapping func(v T) Observable[R]) Operator[T, R] {
	return NewOperator(
		func(source Observable[T]) Observable[R] {
			return switchMapObservable[T, R]{source, mapping}.Subscribe
		},
	)
}

type switchMapObservable[T, R any] struct {
	Source  Observable[T]
	Mapping func(T) Observable[R]
}

func (obs switchMapObservable[T, R]) Subscribe(c Context, sink Observer[R]) {
	c, cancel := c.WithCancel()
	sink = sink.OnTermination(cancel)

	var x struct {
		Context  atomic.Value
		Complete atomic.Bool
		Worker   struct {
			sync.WaitGroup
			Cancel CancelFunc
		}
	}

	x.Context.Store(c.Context)

	startWorker := func(v T) {
		obs1 := obs.Mapping(v)
		w, cancelw := c.WithCancel()

		x.Context.Store(w.Context)
		x.Worker.Add(1)
		x.Worker.Cancel = cancelw

		obs1.Subscribe(w, func(n Notification[R]) {
			switch n.Kind {
			case KindNext:
				sink(n)

			case KindError, KindComplete:
				defer x.Worker.Done()

				cancelw()

				switch n.Kind {
				case KindError:
					if x.Context.CompareAndSwap(w.Context, sentinel) {
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
			if x.Context.Swap(c.Context) == sentinel {
				x.Context.Store(sentinel)
				return
			}

			if x.Worker.Cancel != nil {
				x.Worker.Cancel()
				x.Worker.Wait()
			}

			startWorker(n.Value)

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
