package rx

import (
	"context"
)

type sampleOperator struct {
	Notifier Observable
}

func (op sampleOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	type X struct {
		LatestValue    interface{}
		HasLatestValue bool
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	op.Notifier.Subscribe(ctx, func(t Notification) {
		if x, ok := <-cx; ok {
			if t.HasError {
				close(cx)
				sink(t)
				return
			}
			if x.HasLatestValue {
				sink.Next(x.LatestValue)
				x.HasLatestValue = false
			}
			cx <- x
		}
	})

	if ctx.Err() != nil {
		return Done()
	}

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

// Sample creates an Observable that emits the most recently emitted value from
// the source Observable whenever another Observable, the notifier, emits.
//
// It's like SampleTime, but samples whenever the notifier Observable emits
// something.
func (Operators) Sample(notifier Observable) OperatorFunc {
	return func(source Observable) Observable {
		op := sampleOperator{notifier}
		return source.Lift(op.Call)
	}
}
