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

	var (
		latestValue    interface{}
		hasLatestValue bool
		try            cancellableLocker
	)

	op.Notifier.Subscribe(ctx, func(t Notification) {
		if try.Lock() {
			if t.HasError {
				try.CancelAndUnlock()
				sink(t)
				return
			}
			defer try.Unlock()
			if hasLatestValue {
				sink.Next(latestValue)
				hasLatestValue = false
			}
		}
	})

	select {
	case <-ctx.Done():
		return canceledCtx, nothingToDo
	default:
	}

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
