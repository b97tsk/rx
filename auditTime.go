package rx

import (
	"context"
	"time"
)

type auditTimeOperator struct {
	Duration time.Duration
}

func (op auditTimeOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	type X struct {
		LatestValue interface{}
		Scheduled   bool
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	source.Subscribe(ctx, func(t Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				x.LatestValue = t.Value
				shouldSchedule := !x.Scheduled
				x.Scheduled = true

				cx <- x

				if shouldSchedule {
					scheduleOnce(ctx, op.Duration, func() {
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

	return ctx, cancel
}

// AuditTime ignores source values for a duration, then emits the most recent
// value from the source Observable, then repeats this process.
//
// When it sees a source values, it ignores that plus the next ones for a
// duration, and then it emits the most recent value from the source.
func (Operators) AuditTime(duration time.Duration) OperatorFunc {
	return func(source Observable) Observable {
		op := auditTimeOperator{duration}
		return source.Lift(op.Call)
	}
}
