package rx

import (
	"context"
	"time"
)

type sampleTimeOperator struct {
	Duration time.Duration
}

func (op sampleTimeOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	var (
		latestValue    interface{}
		hasLatestValue bool
		try            cancellableLocker
	)

	schedule(ctx, op.Duration, func() {
		if try.Lock() {
			if hasLatestValue {
				sink.Next(latestValue)
				hasLatestValue = false
			}
			try.Unlock()
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
				sink(t)
			}
		}
	})

	return ctx, cancel
}

// SampleTime creates an Observable that emits the most recently emitted value
// from the source Observable within periodic time intervals.
func (Operators) SampleTime(interval time.Duration) OperatorFunc {
	return func(source Observable) Observable {
		op := sampleTimeOperator{interval}
		return source.Lift(op.Call)
	}
}
