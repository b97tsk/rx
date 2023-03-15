package rx

import (
	"context"
	"sync"
	"sync/atomic"
)

// ExhaustAll flattens a higher-order Observable into a first-order Observable
// by dropping inner Observables while the previous inner Observable has not
// yet completed.
func ExhaustAll[_ Observable[T], T any]() Operator[Observable[T], T] {
	return exhaustMap(identity[Observable[T]])
}

// ExhaustMap converts the source Observable into a higher-order Observable,
// by projecting each source value to an Observable, then flattens it into
// a first-order Observable using ExhaustAll.
func ExhaustMap[T, R any](proj func(v T) Observable[R]) Operator[T, R] {
	if proj == nil {
		panic("proj == nil")
	}

	return exhaustMap(proj)
}

// ExhaustMapTo converts the source Observable into a higher-order Observable,
// by projecting each source value to the same Observable, then flattens it
// into a first-order Observable using ExhaustAll.
func ExhaustMapTo[T, R any](inner Observable[R]) Operator[T, R] {
	return exhaustMap(func(T) Observable[R] { return inner })
}

func exhaustMap[T, R any](proj func(v T) Observable[R]) Operator[T, R] {
	return NewOperator(
		func(source Observable[T]) Observable[R] {
			return exhaustMapObservable[T, R]{source, proj}.Subscribe
		},
	)
}

type exhaustMapObservable[T, R any] struct {
	Source  Observable[T]
	Project func(T) Observable[R]
}

func (obs exhaustMapObservable[T, R]) Subscribe(ctx context.Context, sink Observer[R]) {
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
			if x.Context.Load() == source {
				startWorker(n.Value)
			}

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
