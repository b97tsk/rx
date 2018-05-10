package rx

import (
	"context"
	"time"
)

type auditTimeOperator struct {
	source    Operator
	duration  time.Duration
	scheduler Scheduler
}

func (op auditTimeOperator) ApplyOptions(options []Option) Operator {
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

func (op auditTimeOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	scheduleCtx := canceledCtx
	scheduleDone := scheduleCtx.Done()
	try := cancellableLocker{}
	latestValue := interface{}(nil)

	doSchedule := func() {
		select {
		case <-scheduleDone:
		default:
			return
		}

		scheduleCtx, _ = op.scheduler.ScheduleOnce(ctx, op.duration, func() {
			if try.Lock() {
				defer try.Unlock()
				ob.Next(latestValue)
			}
		})
		scheduleDone = scheduleCtx.Done()
	}

	op.source.Call(ctx, ObserverFunc(func(t Notification) {
		if try.Lock() {
			switch {
			case t.HasValue:
				latestValue = t.Value
				try.Unlock()
				doSchedule()
			case t.HasError:
				try.CancelAndUnlock()
				ob.Error(t.Value.(error))
				cancel()
			default:
				try.CancelAndUnlock()
				ob.Complete()
				cancel()
			}
		}
	}))

	return ctx, cancel
}

// AuditTime ignores source values for a duration, then emits the most recent
// value from the source Observable, then repeats this process.
//
// When it sees a source values, it ignores that plus the next ones for a
// duration, and then it emits the most recent value from the source.
func (o Observable) AuditTime(duration time.Duration) Observable {
	op := auditTimeOperator{
		source:    o.Op,
		duration:  duration,
		scheduler: DefaultScheduler,
	}
	return Observable{op}
}
