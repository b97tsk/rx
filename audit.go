package rx

import (
	"context"
)

type auditObservable struct {
	Source           Observable
	DurationSelector func(interface{}) Observable
}

func (obs auditObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	type X struct {
		LatestValue interface{}
		Scheduled   bool
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	obs.Source.Subscribe(ctx, func(t Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				x.LatestValue = t.Value
				shouldSchedule := !x.Scheduled
				x.Scheduled = true

				cx <- x

				if shouldSchedule {
					scheduleCtx, scheduleCancel := context.WithCancel(ctx)

					var observer Observer
					observer = func(t Notification) {
						observer = NopObserver
						scheduleCancel()
						if x, ok := <-cx; ok {
							if t.HasError {
								close(cx)
								sink(t)
								return
							}
							sink.Next(x.LatestValue)
							x.Scheduled = false
							cx <- x
						}
					}

					obs := obs.DurationSelector(t.Value)
					obs.Subscribe(scheduleCtx, observer.Notify)
				}

			default:
				close(cx)
				sink(t)
			}
		}
	})

	return ctx, cancel
}

// Audit ignores source values for a duration determined by another Observable,
// then emits the most recent value from the source Observable, then repeats
// this process.
//
// It's like AuditTime, but the silencing duration is determined by a second
// Observable.
func (Operators) Audit(durationSelector func(interface{}) Observable) OperatorFunc {
	return func(source Observable) Observable {
		return auditObservable{source, durationSelector}.Subscribe
	}
}
