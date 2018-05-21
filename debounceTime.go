package rx

import (
	"context"
	"time"
)

type debounceTimeOperator struct {
	Duration time.Duration
}

func (op debounceTimeOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	scheduleCancel := doNothing

	var (
		latestValue interface{}
		try         cancellableLocker
	)

	doSchedule := func() {
		scheduleCancel()

		_, scheduleCancel = scheduleOnce(ctx, op.Duration, func() {
			if try.Lock() {
				defer try.Unlock()
				sink.Next(latestValue)
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
			default:
				try.CancelAndUnlock()
				sink(t)
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
	op := debounceTimeOperator{duration}
	return o.Lift(op.Call)
}
