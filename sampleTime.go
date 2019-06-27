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

	type X struct {
		LatestValue    interface{}
		HasLatestValue bool
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	schedule(ctx, op.Duration, func() {
		if x, ok := <-cx; ok {
			if x.HasLatestValue {
				sink.Next(x.LatestValue)
				x.HasLatestValue = false
			}
			cx <- x
		}
	})

	source.Subscribe(ctx, func(t Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				x.LatestValue = t.Value
				x.HasLatestValue = true
				cx <- x
			default:
				close(cx)
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
