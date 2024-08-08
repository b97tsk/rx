package rx

import (
	"sync"
	"sync/atomic"
)

// ExhaustAll flattens a higher-order [Observable] into a first-order
// [Observable] by dropping inner Observables while the previous inner
// [Observable] has not yet completed.
func ExhaustAll[_ Observable[T], T any]() Operator[Observable[T], T] {
	return ExhaustMap(identity[Observable[T]])
}

// ExhaustMapTo converts the source [Observable] into a higher-order
// [Observable], by mapping each source value to the same [Observable],
// then flattens it into a first-order [Observable] using [ExhaustAll].
func ExhaustMapTo[T, R any](inner Observable[R]) Operator[T, R] {
	return ExhaustMap(func(T) Observable[R] { return inner })
}

// ExhaustMap converts the source [Observable] into a higher-order
// [Observable], by mapping each source value to an [Observable], then
// flattens it into a first-order [Observable] using [ExhaustAll].
func ExhaustMap[T, R any](mapping func(v T) Observable[R]) Operator[T, R] {
	return NewOperator(
		func(source Observable[T]) Observable[R] {
			return exhaustMapObservable[T, R]{source, mapping}.Subscribe
		},
	)
}

type exhaustMapObservable[T, R any] struct {
	source  Observable[T]
	mapping func(T) Observable[R]
}

func (ob exhaustMapObservable[T, R]) Subscribe(c Context, o Observer[R]) {
	c, cancel := c.WithCancel()
	o = o.DoOnTermination(cancel)

	var x struct {
		context  atomic.Value
		complete atomic.Bool
		worker   struct {
			sync.WaitGroup
		}
	}

	x.context.Store(c.Context)

	startWorker := func(v T) {
		obs := ob.mapping(v)
		w, cancelw := c.WithCancel()

		x.context.Store(w.Context)
		x.worker.Add(1)

		obs.Subscribe(w, func(n Notification[R]) {
			switch n.Kind {
			case KindNext:
				o.Emit(n)

			case KindComplete, KindError, KindStop:
				defer x.worker.Done()

				cancelw()

				switch n.Kind {
				case KindComplete:
					if x.context.CompareAndSwap(w.Context, c.Context) && x.complete.Load() && x.context.CompareAndSwap(c.Context, sentinel) {
						o.Emit(n)
					}
				case KindError, KindStop:
					if x.context.Swap(sentinel) != sentinel {
						o.Emit(n)
					}
				}
			}
		})
	}

	ob.source.Subscribe(c, func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
			if x.context.Load() == c.Context {
				startWorker(n.Value)
			}

		case KindComplete:
			x.complete.Store(true)

			if x.context.CompareAndSwap(c.Context, sentinel) {
				o.Complete()
			}

		case KindError, KindStop:
			old := x.context.Swap(sentinel)

			cancel()
			x.worker.Wait()

			if old != sentinel {
				switch n.Kind {
				case KindError:
					o.Error(n.Error)
				case KindStop:
					o.Stop(n.Error)
				}
			}
		}
	})
}
