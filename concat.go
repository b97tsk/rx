package rx

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/b97tsk/rx/internal/queue"
)

// Concat creates an Observable that concatenates multiple Observables
// together by sequentially emitting their values, one Observable after
// the other.
func Concat[T any](some ...Observable[T]) Observable[T] {
	if len(some) == 0 {
		return Empty[T]()
	}

	return observables[T](some).Concat
}

// ConcatWith applies [Concat] to the source Observable along with some other
// Observables to create a first-order Observable.
func ConcatWith[T any](some ...Observable[T]) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return observables[T](append([]Observable[T]{source}, some...)).Concat
		},
	)
}

func (some observables[T]) Concat(ctx context.Context, sink Observer[T]) {
	var observer Observer[T]

	done := ctx.Done()

	subscribeToNext := resistReentry(func() {
		select {
		default:
		case <-done:
			sink.Error(ctx.Err())
			return
		}

		if len(some) == 0 {
			sink.Complete()
			return
		}

		obs := some[0]
		some = some[1:]

		obs.Subscribe(ctx, observer)
	})

	observer = func(n Notification[T]) {
		if n.HasValue || n.HasError {
			sink(n)
			return
		}

		subscribeToNext()
	}

	subscribeToNext()
}

// ConcatAll flattens a higher-order Observable into a first-order Observable
// by concatenating the inner Observables in order.
func ConcatAll[_ Observable[T], T any]() Operator[Observable[T], T] {
	return concatMap(identity[Observable[T]])
}

// ConcatMap converts the source Observable into a higher-order Observable,
// by projecting each source value to an Observable, then flattens it into
// a first-order Observable using ConcatAll.
func ConcatMap[T, R any](proj func(v T) Observable[R]) Operator[T, R] {
	if proj == nil {
		panic("proj == nil")
	}

	return concatMap(proj)
}

// ConcatMapTo converts the source Observable into a higher-order Observable,
// by projecting each source value to the same Observable, then flattens it
// into a first-order Observable using ConcatAll.
func ConcatMapTo[T, R any](inner Observable[R]) Operator[T, R] {
	return concatMap(func(T) Observable[R] { return inner })
}

func concatMap[T, R any](proj func(v T) Observable[R]) Operator[T, R] {
	return NewOperator(
		func(source Observable[T]) Observable[R] {
			return concatMapObservable[T, R]{source, proj}.Subscribe
		},
	)
}

type concatMapObservable[T, R any] struct {
	Source  Observable[T]
	Project func(T) Observable[R]
}

func (obs concatMapObservable[T, R]) Subscribe(ctx context.Context, sink Observer[R]) {
	source, cancelSource := context.WithCancel(ctx)

	sink = sink.OnLastNotification(cancelSource)

	var x struct {
		Context  atomic.Value
		Complete atomic.Bool
		Queue    struct {
			sync.Mutex
			queue.Queue[T]
		}
		Worker struct {
			sync.WaitGroup
		}
	}

	x.Context.Store(source)

	var startWorker func()

	startWorker = resistReentry(func() {
		v := x.Queue.Pop()

		x.Queue.Unlock()

		worker, cancelWorker := context.WithCancel(source)
		x.Context.Store(worker)

		x.Worker.Add(1)

		obs.Project(v).Subscribe(worker, func(n Notification[R]) {
			switch {
			case n.HasValue:
				sink(n)
				return

			default:
				x.Queue.Lock()

				if !n.HasError {
					select {
					default:
					case <-worker.Done():
						n = Error[R](worker.Err())
					}
				}

				if n.HasError {
					swapped := x.Context.CompareAndSwap(worker, sentinel)

					x.Queue.Init()
					x.Queue.Unlock()

					if swapped {
						sink(n)
					}

					break
				}

				if x.Queue.Len() == 0 {
					swapped := x.Context.CompareAndSwap(worker, source)

					x.Queue.Unlock()

					if swapped && x.Complete.Load() && x.Context.CompareAndSwap(source, sentinel) {
						sink.Complete()
					}

					break
				}

				startWorker()
			}

			cancelWorker()
			x.Worker.Done()
		})
	})

	obs.Source.Subscribe(source, func(n Notification[T]) {
		switch {
		case n.HasValue:
			x.Queue.Lock()

			ctx := x.Context.Load()
			if ctx != sentinel {
				x.Queue.Push(n.Value)
			}

			if ctx == source {
				startWorker()
				return
			}

			x.Queue.Unlock()

		case n.HasError:
			ctx := x.Context.Swap(source)

			cancelSource()
			x.Worker.Wait()

			if x.Context.Swap(sentinel) != sentinel && ctx != sentinel {
				sink.Error(n.Error)
			}

			x.Queue.Init()

		default:
			x.Complete.Store(true)

			if x.Context.CompareAndSwap(source, sentinel) {
				sink.Complete()
			}
		}
	})
}
