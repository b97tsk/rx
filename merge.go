package rx

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/b97tsk/rx/internal/queue"
	"github.com/b97tsk/rx/internal/waitgroup"
)

// Merge creates an Observable that concurrently emits all values from every
// given input Observable.
func Merge[T any](some ...Observable[T]) Observable[T] {
	if len(some) == 0 {
		return Empty[T]()
	}

	return observables[T](some).Merge
}

// MergeWith applies [Merge] to the source Observable along with some other
// Observables to create a first-order Observable.
func MergeWith[T any](some ...Observable[T]) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return observables[T](append([]Observable[T]{source}, some...)).Merge
		},
	)
}

func (some observables[T]) Merge(ctx context.Context, sink Observer[T]) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.OnLastNotification(cancel).Serialized()

	var workers atomic.Uint32

	workers.Store(uint32(len(some)))

	observer := func(n Notification[T]) {
		if n.HasValue || n.HasError || workers.Add(^uint32(0)) == 0 {
			sink(n)
		}
	}

	ctxHoisted := waitgroup.Hoist(ctx)

	for _, obs := range some {
		obs := obs

		Go(ctxHoisted, func() { obs.Subscribe(ctx, observer) })
	}
}

// MergeAll flattens a higher-order Observable into a first-order Observable
// which concurrently delivers all values that are emitted on the inner
// Observables.
func MergeAll[_ Observable[T], T any]() MergeMapOperator[Observable[T], T] {
	return MergeMapOperator[Observable[T], T]{
		opts: mergeMapConfig[Observable[T], T]{
			Project:     identity[Observable[T]],
			Concurrency: -1,
		},
	}
}

// MergeMap converts the source Observable into a higher-order Observable,
// by projecting each source value to an Observable, then flattens it into
// a first-order Observable using MergeAll.
func MergeMap[T, R any](proj func(v T) Observable[R]) MergeMapOperator[T, R] {
	if proj == nil {
		panic("proj == nil")
	}

	return MergeMapOperator[T, R]{
		opts: mergeMapConfig[T, R]{
			Project:     proj,
			Concurrency: -1,
		},
	}
}

// MergeMapTo converts the source Observable into a higher-order Observable,
// by projecting each source value to the same Observable, then flattens it
// into a first-order Observable using MergeAll.
func MergeMapTo[T, R any](inner Observable[R]) MergeMapOperator[T, R] {
	return MergeMapOperator[T, R]{
		opts: mergeMapConfig[T, R]{
			Project:     func(T) Observable[R] { return inner },
			Concurrency: -1,
		},
	}
}

type mergeMapConfig[T, R any] struct {
	Project     func(T) Observable[R]
	Concurrency int
}

// MergeMapOperator is an [Operator] type for [MergeMap].
type MergeMapOperator[T, R any] struct {
	opts mergeMapConfig[T, R]
}

// WithConcurrency sets Concurrency option to a given value.
// It must not be zero. The default value is -1 (unlimited).
func (op MergeMapOperator[T, R]) WithConcurrency(n int) MergeMapOperator[T, R] {
	if n == 0 {
		panic("n == 0")
	}

	op.opts.Concurrency = n

	return op
}

// Apply implements the Operator interface.
func (op MergeMapOperator[T, R]) Apply(source Observable[T]) Observable[R] {
	return mergeMapObservable[T, R]{source, op.opts}.Subscribe
}

type mergeMapObservable[T, R any] struct {
	Source Observable[T]
	mergeMapConfig[T, R]
}

func (obs mergeMapObservable[T, R]) Subscribe(ctx context.Context, sink Observer[R]) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.OnLastNotification(cancel).Serialized()

	var x struct {
		Workers  atomic.Uint32
		Complete atomic.Bool
		Queue    struct {
			sync.Mutex
			queue.Queue[T]
			Sealed bool
		}
	}

	var startWorker func()

	ctxHoisted := waitgroup.Hoist(ctx)

	startWorker = func() {
		obs1 := obs.Project(x.Queue.Pop())

		x.Queue.Unlock()

		Go(ctxHoisted, func() {
			obs1.Subscribe(ctx, func(n Notification[R]) {
				if n.HasValue {
					sink(n)
					return
				}

				x.Queue.Lock()

				if n.HasError {
					x.Queue.Sealed = true

					x.Queue.Init()
					x.Queue.Unlock()

					sink(n)

					return
				}

				if x.Queue.Len() == 0 {
					workers := x.Workers.Add(^uint32(0))

					x.Queue.Unlock()

					if workers == 0 && x.Complete.Load() && x.Workers.Load() == 0 {
						sink(n)
					}

					return
				}

				startWorker()
			})
		})
	}

	obs.Source.Subscribe(ctx, func(n Notification[T]) {
		switch {
		case n.HasValue:
			x.Queue.Lock()

			if x.Queue.Sealed {
				x.Queue.Unlock()
				return
			}

			x.Queue.Push(n.Value)

			if x.Workers.Load() != uint32(obs.Concurrency) {
				x.Workers.Add(1)

				startWorker()

				return
			}

			x.Queue.Unlock()

		case n.HasError:
			sink.Error(n.Error)

		default:
			x.Complete.Store(true)

			if x.Workers.Load() == 0 {
				sink.Complete()
			}
		}
	})
}
