package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

type auditObservable struct {
	Source           rx.Observable
	DurationSelector func(interface{}) rx.Observable
}

func (obs auditObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	type X struct {
		LatestValue interface{}
		Scheduled   bool
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				x.LatestValue = t.Value
				shouldSchedule := !x.Scheduled
				x.Scheduled = true

				cx <- x

				if shouldSchedule {
					scheduleCtx, scheduleCancel := context.WithCancel(ctx)

					var observer rx.Observer
					observer = func(t rx.Notification) {
						observer = rx.NopObserver
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
}

// Audit ignores source values for a duration determined by another Observable,
// then emits the most recent value from the source Observable, then repeats
// this process.
//
// It's like AuditTime, but the silencing duration is determined by a second
// Observable.
func Audit(durationSelector func(interface{}) rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := auditObservable{source, durationSelector}
		return rx.Create(obs.Subscribe)
	}
}