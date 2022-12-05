package rx

import (
	"context"
	"time"

	"github.com/b97tsk/rx/internal/critical"
)

// Debounce emits a value from the source Observable only after a particular
// time span, determined by another Observable, has passed without another
// source emission.
//
// It's like DebounceTime, but the time span of emission silence is determined
// by a second Observable.
func Debounce[T, U any](durationSelector func(v T) Observable[U]) Operator[T, T] {
	if durationSelector == nil {
		panic("durationSelector == nil")
	}

	return debounce(durationSelector)
}

// DebounceTime emits a value from the source Observable only after
// a particular time span has passed without another source emission.
func DebounceTime[T any](d time.Duration) Operator[T, T] {
	obsTimer := Timer(d)

	durationSelector := func(T) Observable[time.Time] { return obsTimer }

	return debounce(durationSelector)
}

func debounce[T, U any](durationSelector func(v T) Observable[U]) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return debounceObservable[T, U]{source, durationSelector}.Subscribe
		},
	)
}

type debounceObservable[T, U any] struct {
	Source           Observable[T]
	DurationSelector func(T) Observable[U]
}

func (obs debounceObservable[T, U]) Subscribe(ctx context.Context, sink Observer[T]) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel)

	var x struct {
		critical.Section
		Latest struct {
			Value    T
			HasValue bool
		}
	}

	var cancelSchedule context.CancelFunc

	obs.Source.Subscribe(ctx, func(n Notification[T]) {
		if critical.Enter(&x.Section) {
			switch {
			case n.HasValue:
				x.Latest.Value = n.Value
				x.Latest.HasValue = true

				critical.Leave(&x.Section)

				if cancelSchedule != nil {
					cancelSchedule()
				}

				ctx, cancel := context.WithCancel(ctx)
				cancelSchedule = cancel

				var noop bool

				obs.DurationSelector(n.Value).Subscribe(ctx, func(n Notification[U]) {
					if noop {
						return
					}

					noop = true

					if ctx.Err() != nil {
						return
					}

					cancel()

					if critical.Enter(&x.Section) {
						switch {
						case n.HasValue:
							if x.Latest.HasValue {
								sink.Next(x.Latest.Value)
								x.Latest.HasValue = false
							}

							critical.Leave(&x.Section)

						case n.HasError:
							critical.Close(&x.Section)

							sink.Error(n.Error)

						default:
							critical.Leave(&x.Section)
						}
					}
				})

			case n.HasError:
				critical.Close(&x.Section)

				sink(n)

			default:
				critical.Close(&x.Section)

				if x.Latest.HasValue {
					sink.Next(x.Latest.Value)
				}

				sink(n)
			}
		}
	})
}
