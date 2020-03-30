package operators

import (
	"context"
	"time"

	"github.com/b97tsk/rx"
)

type auditTimeObservable struct {
	Source   rx.Observable
	Duration time.Duration
}

func (obs auditTimeObservable) Subscribe(ctx context.Context, sink rx.Observer) {
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
					scheduleOnce(ctx, obs.Duration, func() {
						if x, ok := <-cx; ok {
							sink.Next(x.LatestValue)
							x.Scheduled = false
							cx <- x
						}
					})
				}

			default:
				close(cx)
				sink(t)
			}
		}
	})
}

// AuditTime ignores source values for a duration, then emits the most recent
// value from the source Observable, then repeats this process.
//
// When it sees a source values, it ignores that plus the next ones for a
// duration, and then it emits the most recent value from the source.
func AuditTime(duration time.Duration) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := auditTimeObservable{source, duration}
		return rx.Create(obs.Subscribe)
	}
}
