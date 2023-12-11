package rx

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// Throttle emits a value from the source Observable, then ignores
// subsequent source values for a duration determined by another Observable,
// then repeats this process until the source completes.
//
// It's like [ThrottleTime], but the silencing duration is determined
// by a second Observable.
func Throttle[T, U any](durationSelector func(v T) Observable[U]) ThrottleOperator[T, U] {
	if durationSelector == nil {
		panic("durationSelector == nil")
	}

	return ThrottleOperator[T, U]{
		opts: throttleConfig[T, U]{
			DurationSelector: durationSelector,
			Leading:          true,
			Trailing:         false,
		},
	}
}

// ThrottleTime emits a value from the source Observable, then ignores
// subsequent source values for a duration, then repeats this process until
// the source completes.
//
// ThrottleTime lets a value pass, then ignores source values
// for the next duration time.
func ThrottleTime[T any](d time.Duration) ThrottleOperator[T, time.Time] {
	obsTimer := Timer(d)

	durationSelector := func(T) Observable[time.Time] { return obsTimer }

	return ThrottleOperator[T, time.Time]{
		opts: throttleConfig[T, time.Time]{
			DurationSelector: durationSelector,
			Leading:          true,
			Trailing:         false,
		},
	}
}

type throttleConfig[T, U any] struct {
	DurationSelector func(T) Observable[U]
	Leading          bool
	Trailing         bool
}

// ThrottleOperator is an [Operator] type for [Throttle].
type ThrottleOperator[T, U any] struct {
	opts throttleConfig[T, U]
}

// WithLeading sets Leading option to a given value.
func (op ThrottleOperator[T, U]) WithLeading(v bool) ThrottleOperator[T, U] {
	op.opts.Leading = v
	return op
}

// WithTrailing sets Trailing option to a given value.
func (op ThrottleOperator[T, U]) WithTrailing(v bool) ThrottleOperator[T, U] {
	op.opts.Trailing = v
	return op
}

// Apply implements the Operator interface.
func (op ThrottleOperator[T, U]) Apply(source Observable[T]) Observable[T] {
	return throttleObservable[T, U]{source, op.opts}.Subscribe
}

type throttleObservable[T, U any] struct {
	Source Observable[T]
	throttleConfig[T, U]
}

func (obs throttleObservable[T, U]) Subscribe(ctx context.Context, sink Observer[T]) {
	source, cancelSource := context.WithCancel(ctx)

	sink = sink.OnLastNotification(cancelSource)

	var x struct {
		Context  atomic.Value
		Complete atomic.Bool
		Trailing struct {
			sync.Mutex
			Value    T
			HasValue atomic.Bool
		}
		Worker struct {
			sync.WaitGroup
		}
	}

	x.Context.Store(source)

	var doThrottle func(T)

	doThrottle = func(v T) {
		worker, cancelWorker := context.WithCancel(source)

		x.Context.Store(worker)

		x.Worker.Add(1)

		var noop bool

		obs.DurationSelector(v).Subscribe(worker, func(n Notification[U]) {
			if noop {
				return
			}

			noop = true

			cancelWorker()

			switch n.Kind {
			case KindNext:
				if obs.Trailing && x.Trailing.HasValue.Load() {
					x.Trailing.Lock()
					value := x.Trailing.Value
					x.Trailing.HasValue.Store(false)
					x.Trailing.Unlock()
					sink.Next(value)

					if !x.Complete.Load() {
						doThrottle(value)
					}
				}

				if x.Context.CompareAndSwap(worker, source) && x.Complete.Load() && x.Context.CompareAndSwap(source, sentinel) {
					sink.Complete()
				}

			case KindError:
				if x.Context.CompareAndSwap(worker, sentinel) {
					sink.Error(n.Error)
				}

			case KindComplete:
				if x.Context.CompareAndSwap(worker, source) && x.Complete.Load() && x.Context.CompareAndSwap(source, sentinel) {
					sink.Complete()
				}
			}

			x.Worker.Done()
		})
	}

	obs.Source.Subscribe(source, func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
			x.Trailing.Lock()
			x.Trailing.Value = n.Value
			x.Trailing.HasValue.Store(true)
			x.Trailing.Unlock()

			if x.Context.Load() == source {
				if obs.Leading {
					x.Trailing.HasValue.Store(false)
					sink(n)
				}

				doThrottle(n.Value)
			}

		case KindError:
			old := x.Context.Swap(sentinel)

			cancelSource()
			x.Worker.Wait()

			if old != sentinel {
				sink(n)
			}

		case KindComplete:
			x.Complete.Store(true)

			if x.Context.CompareAndSwap(source, sentinel) {
				sink(n)
			}
		}
	})
}
