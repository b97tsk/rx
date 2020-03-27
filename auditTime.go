package rx

import (
	"context"
	"time"
)

type auditTimeObservable struct {
	Source   Observable
	Duration time.Duration
}

func (obs auditTimeObservable) Subscribe(ctx context.Context, sink Observer) {
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
func (Operators) AuditTime(duration time.Duration) Operator {
	return func(source Observable) Observable {
		obs := auditTimeObservable{source, duration}
		return Create(obs.Subscribe)
	}
}
