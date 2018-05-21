package rx

import (
	"context"
	"time"
)

type throttleTimeOperator struct {
	Duration time.Duration
}

func (op throttleTimeOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	scheduleCtx := canceledCtx
	scheduleDone := scheduleCtx.Done()

	source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			select {
			case <-scheduleDone:
			default:
				return
			}

			sink(t)

			scheduleCtx, _ = scheduleOnce(ctx, op.Duration, doNothing)
			scheduleDone = scheduleCtx.Done()

		default:
			sink(t)
			cancel()
		}
	})

	return ctx, cancel
}

// ThrottleTime creates an Observable that emits a value from the source
// Observable, then ignores subsequent source values for a duration, then
// repeats this process.
//
// ThrottleTime lets a value pass, then ignores source values for the next
// duration time.
func (o Observable) ThrottleTime(duration time.Duration) Observable {
	op := throttleTimeOperator{duration}
	return o.Lift(op.Call)
}
