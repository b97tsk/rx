package rx

import (
	"context"
	"time"
)

type auditTimeOperator struct {
	duration  time.Duration
	scheduler Scheduler
}

func (op auditTimeOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	scheduleCtx := canceledCtx
	scheduleDone := scheduleCtx.Done()

	var (
		latestValue interface{}
		try         cancellableLocker
	)

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

	source.Subscribe(ctx, func(t Notification) {
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
	})

	return ctx, cancel
}

// AuditTime ignores source values for a duration, then emits the most recent
// value from the source Observable, then repeats this process.
//
// When it sees a source values, it ignores that plus the next ones for a
// duration, and then it emits the most recent value from the source.
func (o Observable) AuditTime(duration time.Duration) Observable {
	op := auditTimeOperator{duration, DefaultScheduler}
	return o.Lift(op.Call)
}
