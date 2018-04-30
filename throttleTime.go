package rx

import (
	"context"
	"time"
)

type throttleTimeOperator struct {
	source    Operator
	duration  time.Duration
	scheduler Scheduler
}

func (op throttleTimeOperator) ApplyOptions(options []Option) Operator {
	for _, opt := range options {
		switch t := opt.(type) {
		case schedulerOption:
			op.scheduler = t.Value
		default:
			panic(ErrUnsupportedOption)
		}
	}
	return op
}

func (op throttleTimeOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	scheduleCtx := canceledCtx
	scheduleDone := scheduleCtx.Done()

	op.source.Call(ctx, ObserverFunc(func(t Notification) {
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
	}))

	return ctx, cancel
}

// ThrottleTime creates an Observable that emits a value from the source
// Observable, then ignores subsequent source values for a duration, then
// repeats this process.
//
// ThrottleTime lets a value pass, then ignores source values for the next
// duration time.
func (o Observable) ThrottleTime(duration time.Duration) Observable {
	op := throttleTimeOperator{
		source:    o.Op,
		duration:  duration,
		scheduler: DefaultScheduler,
	}
	return Observable{op}
}
