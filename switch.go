package rx

import (
	"context"
	"sync"
	"sync/atomic"
)

// SwitchAll flattens a higher-order Observable into a first-order Observable
// by subscribing to only the most recently emitted of those inner Observables.
func SwitchAll[_ Observable[T], T any]() Operator[Observable[T], T] {
	return switchMap(identity[Observable[T]])
}

// SwitchMap converts the source Observable into a higher-order Observable,
// by projecting each source value to an Observable, then flattens it into
// a first-order Observable using SwitchAll.
func SwitchMap[T, R any](proj func(v T) Observable[R]) Operator[T, R] {
	if proj == nil {
		panic("proj == nil")
	}

	return switchMap(proj)
}

// SwitchMapTo converts the source Observable into a higher-order Observable,
// by projecting each source value to the same Observable, then flattens it
// into a first-order Observable using SwitchAll.
func SwitchMapTo[T, R any](inner Observable[R]) Operator[T, R] {
	return switchMap(func(T) Observable[R] { return inner })
}

func switchMap[T, R any](proj func(v T) Observable[R]) Operator[T, R] {
	return NewOperator(
		func(source Observable[T]) Observable[R] {
			return switchMapObservable[T, R]{source, proj}.Subscribe
		},
	)
}

type switchMapObservable[T, R any] struct {
	Source  Observable[T]
	Project func(T) Observable[R]
}

func (obs switchMapObservable[T, R]) Subscribe(ctx context.Context, sink Observer[R]) {
	source, cancel := context.WithCancel(ctx)

	sink = sink.OnLastNotification(cancel)

	var x struct {
		Context  atomic.Value
		Complete atomic.Bool
		Worker   struct {
			sync.WaitGroup
			Cancel context.CancelFunc
		}
	}

	x.Context.Store(source)

	startWorker := func(v T) {
		worker, cancel := context.WithCancel(source)
		x.Context.Store(worker)

		x.Worker.Cancel = cancel

		x.Worker.Add(1)

		obs.Project(v).Subscribe(worker, func(n Notification[R]) {
			switch {
			case n.HasValue:
				sink(n)
				return

			case n.HasError:
				if x.Context.CompareAndSwap(worker, sentinel) {
					sink(n)
				}

			default:
				if x.Context.CompareAndSwap(worker, source) && x.Complete.Load() && x.Context.CompareAndSwap(source, sentinel) {
					sink(n)
				}
			}

			cancel()
			x.Worker.Done()
		})
	}

	obs.Source.Subscribe(source, func(n Notification[T]) {
		switch {
		case n.HasValue:
			if x.Context.Swap(source) == sentinel {
				x.Context.Store(sentinel)
				return
			}

			if x.Worker.Cancel != nil {
				x.Worker.Cancel()
				x.Worker.Wait()
			}

			startWorker(n.Value)

		case n.HasError:
			ctx := x.Context.Swap(source)

			if x.Worker.Cancel != nil {
				x.Worker.Cancel()
				x.Worker.Wait()
			}

			if x.Context.Swap(sentinel) != sentinel && ctx != sentinel {
				sink.Error(n.Error)
			}

		default:
			x.Complete.Store(true)

			if x.Context.CompareAndSwap(source, sentinel) {
				sink.Complete()
			}
		}
	})
}
