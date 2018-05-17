package rx

import (
	"context"
	"time"
)

type debounceTimeOperator struct {
	duration  time.Duration
	scheduler Scheduler
}

func (op debounceTimeOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	scheduleCancel := noopFunc

	var (
		latestValue interface{}
		try         cancellableLocker
	)

	doSchedule := func() {
		scheduleCancel()

		_, scheduleCancel = op.scheduler.ScheduleOnce(ctx, op.duration, func() {
			if try.Lock() {
				defer try.Unlock()
				ob.Next(latestValue)
			}
		})
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

// DebounceTime creates an Observable that emits a value from the source
// Observable only after a particular time span has passed without another
// source emission.
//
// It's like Delay, but passes only the most recent value from each burst of
// emissions.
func (o Observable) DebounceTime(duration time.Duration) Observable {
	op := debounceTimeOperator{duration, DefaultScheduler}
	return o.Lift(op.Call)
}
