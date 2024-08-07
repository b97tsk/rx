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
	source  Observable[T]
	mapping func(T) Observable[R]
}

func (ob switchMapObservable[T, R]) Subscribe(c Context, o Observer[R]) {
	c, cancel := c.WithCancel()
	o = o.DoOnTermination(cancel)

	var x struct {
		context  atomic.Value
		complete atomic.Bool
		worker   struct {
			sync.WaitGroup
			cancel CancelFunc
		}
	}

	x.context.Store(c.Context)

	startWorker := func(v T) {
		obs := ob.mapping(v)
		w, cancelw := c.WithCancel()

		x.context.Store(w.Context)
		x.worker.Add(1)
		x.worker.cancel = cancelw

		obs.Subscribe(w, func(n Notification[R]) {
			switch n.Kind {
			case KindNext:
				o.Emit(n)

			case KindError, KindComplete:
				defer x.worker.Done()

				cancelw()

				switch n.Kind {
				case KindError:
					if x.context.CompareAndSwap(w.Context, sentinel) {
						o.Emit(n)
					}
				case KindComplete:
					if x.context.CompareAndSwap(w.Context, c.Context) && x.complete.Load() && x.context.CompareAndSwap(c.Context, sentinel) {
						o.Emit(n)
					}
				}
			}
		})
	}

	ob.source.Subscribe(c, func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
			if x.context.Swap(c.Context) == sentinel {
				x.context.Store(sentinel)
				return
			}

			if x.worker.cancel != nil {
				x.worker.cancel()
				x.worker.Wait()
			}

			startWorker(n.Value)

		case KindError:
			old := x.context.Swap(sentinel)

			cancel()
			x.worker.Wait()

			if old != sentinel {
				o.Error(n.Error)
			}

		case KindComplete:
			x.complete.Store(true)

			if x.context.CompareAndSwap(c.Context, sentinel) {
				o.Complete()
			}
		}
	})
}
