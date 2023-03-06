package rx

import (
	"context"
	"time"

	"github.com/b97tsk/rx/internal/critical"
)

// Audit ignores source values for a duration determined by another Observable,
// then emits the most recent value from the source Observable, then repeats
// this process.
//
// It's like [AuditTime], but the silencing duration is determined by a second
// Observable.
func Audit[T, U any](durationSelector func(v T) Observable[U]) Operator[T, T] {
	if durationSelector == nil {
		panic("durationSelector == nil")
	}

	return audit(durationSelector)
}

// AuditTime ignores source values for a duration, then emits the most recent
// value from the source Observable, then repeats this process.
//
// When it sees a source value, it ignores that plus the next ones for a
// duration, and then it emits the most recent value from the source.
func AuditTime[T any](d time.Duration) Operator[T, T] {
	obsTimer := Timer(d)

	durationSelector := func(T) Observable[time.Time] { return obsTimer }

	return audit(durationSelector)
}

func audit[T, U any](durationSelector func(v T) Observable[U]) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return auditObservable[T, U]{source, durationSelector}.Subscribe
		},
	)
}

type auditObservable[T, U any] struct {
	Source           Observable[T]
	DurationSelector func(T) Observable[U]
}

func (obs auditObservable[T, U]) Subscribe(ctx context.Context, sink Observer[T]) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.OnLastNotification(cancel)

	var x struct {
		critical.Section
		LatestValue T
		Scheduled   bool
		Complete    bool
	}

	obs.Source.Subscribe(ctx, func(n Notification[T]) {
		if critical.Enter(&x.Section) {
			switch {
			case n.HasValue:
				shouldSchedule := !x.Scheduled

				x.LatestValue = n.Value
				x.Scheduled = true

				critical.Leave(&x.Section)

				if shouldSchedule {
					ctx, cancel := context.WithCancel(ctx)

					var noop bool

					obs.DurationSelector(n.Value).Subscribe(ctx, func(n Notification[U]) {
						if noop {
							return
						}

						noop = true

						cancel()

						if critical.Enter(&x.Section) {
							x.Scheduled = false

							switch {
							case n.HasValue:
								sink.Next(x.LatestValue)

								if x.Complete {
									sink.Complete()
								}

								critical.Leave(&x.Section)

							case n.HasError:
								critical.Close(&x.Section)

								sink.Error(n.Error)

							default:
								if x.Complete {
									sink.Complete()
								}

								critical.Leave(&x.Section)
							}
						}
					})
				}

			case n.HasError:
				critical.Close(&x.Section)

				sink(n)

			default:
				x.Complete = true

				if x.Scheduled {
					critical.Leave(&x.Section)
					break
				}

				critical.Close(&x.Section)

				sink(n)
			}
		}
	})
}
