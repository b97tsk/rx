package rx

import (
	"context"
	"time"

	"github.com/b97tsk/rx/internal/critical"
)

// Throttle emits a value from the source Observable, then ignores subsequent
// source values for a duration determined by another Observable, then repeats
// this process.
//
// It's like ThrottleTime, but the silencing duration is determined by a second
// Observable.
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
// subsequent source values for a duration, then repeats this process.
//
// ThrottleTime lets a value pass, then ignores source values for the next
// duration time.
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

// ThrottleOperator is an Operator type for Throttle.
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

// AsOperator converts op to an Operator.
//
// Once type inference has improved in Go, this method will be removed.
func (op ThrottleOperator[T, U]) AsOperator() Operator[T, T] { return op }

type throttleObservable[T, U any] struct {
	Source Observable[T]
	throttleConfig[T, U]
}

func (obs throttleObservable[T, U]) Subscribe(ctx context.Context, sink Observer[T]) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel)

	var x struct {
		critical.Section
		Trailing struct {
			Value    T
			HasValue bool
		}
		Throttling bool
		Completed  bool
	}

	var doThrottle func(T)

	doThrottle = func(v T) {
		x.Throttling = true

		ctx, cancel := context.WithCancel(ctx)

		var noop bool

		go obs.DurationSelector(v).Subscribe(ctx, func(n Notification[U]) {
			if noop {
				return
			}

			noop = true

			cancel()

			if critical.Enter(&x.Section) {
				x.Throttling = false

				switch {
				case n.HasValue:
					if obs.Trailing && x.Trailing.HasValue {
						sink.Next(x.Trailing.Value)
						x.Trailing.HasValue = false

						if !x.Completed {
							doThrottle(x.Trailing.Value)
						}
					}

					if x.Completed {
						sink.Complete()
					}

					critical.Leave(&x.Section)

				case n.HasError:
					critical.Close(&x.Section)

					sink.Error(n.Error)

				default:
					if x.Completed {
						sink.Complete()
					}

					critical.Leave(&x.Section)
				}
			}
		})
	}

	obs.Source.Subscribe(ctx, func(n Notification[T]) {
		if critical.Enter(&x.Section) {
			switch {
			case n.HasValue:
				x.Trailing.Value = n.Value
				x.Trailing.HasValue = true

				if !x.Throttling {
					doThrottle(n.Value)

					if obs.Leading {
						x.Trailing.HasValue = false

						sink(n)
					}
				}

				critical.Leave(&x.Section)

			case n.HasError:
				critical.Close(&x.Section)

				sink(n)

			default:
				x.Completed = true

				if obs.Trailing && x.Trailing.HasValue && x.Throttling {
					critical.Leave(&x.Section)
					break
				}

				critical.Close(&x.Section)

				sink(n)
			}
		}
	})
}
