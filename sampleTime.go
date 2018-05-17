package rx

import (
	"context"
	"time"
)

type sampleTimeOperator struct {
	interval  time.Duration
	scheduler Scheduler
}

func (op sampleTimeOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	var (
		latestValue    interface{}
		hasLatestValue bool
		try            cancellableLocker
	)

	op.scheduler.Schedule(ctx, op.interval, func() {
		if try.Lock() {
			defer try.Unlock()
			if hasLatestValue {
				ob.Next(latestValue)
				hasLatestValue = false
			}
		}
	})

	source.Subscribe(ctx, func(t Notification) {
		if try.Lock() {
			switch {
			case t.HasValue:
				latestValue = t.Value
				hasLatestValue = true
				try.Unlock()
			default:
				try.CancelAndUnlock()
				t.Observe(ob)
				cancel()
			}
		}
	})

	return ctx, cancel
}

// SampleTime creates an Observable that emits the most recently emitted value
// from the source Observable within periodic time intervals.
func (o Observable) SampleTime(interval time.Duration) Observable {
	op := sampleTimeOperator{interval, DefaultScheduler}
	return o.Lift(op.Call)
}
