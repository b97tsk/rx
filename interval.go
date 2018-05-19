package rx

import (
	"context"
	"time"
)

type intervalOperator struct {
	initialDelay time.Duration
	period       time.Duration
	scheduler    Scheduler
}

func (op intervalOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	if op.initialDelay != op.period {
		op.scheduler.ScheduleOnce(ctx, op.initialDelay, func() {
			index := 0
			wait := make(chan struct{})
			op.scheduler.Schedule(ctx, op.period, func() {
				<-wait
				sink.Next(index)
				index++
			})
			sink.Next(index)
			index++
			close(wait)
		})
	} else {
		index := 0
		op.scheduler.Schedule(ctx, op.period, func() {
			sink.Next(index)
			index++
		})
	}

	return ctx, cancel
}

// Interval creates an Observable that emits sequential integers every
// specified interval of time, on a specified Scheduler.
func Interval(period time.Duration) Observable {
	op := intervalOperator{period, period, DefaultScheduler}
	return Observable{}.Lift(op.Call)
}

// Timer creates an Observable that starts emitting after an initialDelay and
// emits ever increasing integers after each period of time thereafter.
//
// Its like Interval, but you can specify when should the emissions start.
func Timer(initialDelay, period time.Duration) Observable {
	op := intervalOperator{initialDelay, period, DefaultScheduler}
	return Observable{}.Lift(op.Call)
}
