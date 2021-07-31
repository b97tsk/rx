package operators

import (
	"context"
	"time"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/critical"
)

// Audit ignores source values for a duration determined by another Observable,
// then emits the most recent value from the source, then repeats this process.
//
// It's like AuditTime, but the silencing duration is determined by a second
// Observable.
func Audit(durationSelector func(interface{}) rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return auditObservable{source, durationSelector}.Subscribe
	}
}

// AuditTime ignores source values for a duration, then emits the most recent
// value from the source, then repeats this process.
//
// When it sees a source values, it ignores that plus the next ones for a
// duration, and then it emits the most recent value from the source.
func AuditTime(d time.Duration) rx.Operator {
	obsTimer := rx.Timer(d)

	durationSelector := func(interface{}) rx.Observable { return obsTimer }

	return Audit(durationSelector)
}

type auditObservable struct {
	Source           rx.Observable
	DurationSelector func(interface{}) rx.Observable
}

func (obs auditObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel)

	var x struct {
		critical.Section
		LatestValue interface{}
		Scheduled   bool
		Completed   bool
	}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		if critical.Enter(&x.Section) {
			switch {
			case t.HasValue:
				shouldSchedule := !x.Scheduled

				x.LatestValue = t.Value
				x.Scheduled = true

				critical.Leave(&x.Section)

				if shouldSchedule {
					scheduleCtx, scheduleCancel := context.WithCancel(ctx)

					var observer rx.Observer

					observer = func(t rx.Notification) {
						observer = rx.Noop

						scheduleCancel()

						if critical.Enter(&x.Section) {
							x.Scheduled = false

							switch {
							case t.HasValue:
								defer critical.Leave(&x.Section)

								sink.Next(x.LatestValue)

								if x.Completed {
									sink.Complete()
								}

							case t.HasError:
								critical.Close(&x.Section)

								sink(t)

							default:
								defer critical.Leave(&x.Section)

								if x.Completed {
									sink(t)
								}
							}
						}
					}

					obs1 := obs.DurationSelector(t.Value)

					obs1.Subscribe(scheduleCtx, observer.Sink)
				}

			case t.HasError:
				critical.Close(&x.Section)

				sink(t)

			default:
				x.Completed = true

				if x.Scheduled {
					critical.Leave(&x.Section)
					break
				}

				critical.Close(&x.Section)

				sink(t)
			}
		}
	})
}
