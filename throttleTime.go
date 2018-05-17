package rx

import (
	"context"
	"time"
)

type throttleTimeOperator struct {
	duration  time.Duration
	scheduler Scheduler
}

func (op throttleTimeOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
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

			ob.Next(t.Value)

			scheduleCtx, _ = op.scheduler.ScheduleOnce(ctx, op.duration, noopFunc)
			scheduleDone = scheduleCtx.Done()

		case t.HasError:
			ob.Error(t.Value.(error))
			cancel()

		default:
			ob.Complete()
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
	op := throttleTimeOperator{duration, DefaultScheduler}
	return o.Lift(op.Call)
}
