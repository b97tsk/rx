package rx

import (
	"context"
)

type sampleOperator struct {
	notifier Observable
}

func (op sampleOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	var (
		latestValue    interface{}
		hasLatestValue bool
		try            cancellableLocker
	)

	op.notifier.Subscribe(ctx, func(t Notification) {
		if t.HasError {
			t.Observe(ob)
			cancel()
			return
		}
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

// Sample creates an Observable that emits the most recently emitted value from
// the source Observable whenever another Observable, the notifier, emits.
//
// It's like SampleTime, but samples whenever the notifier Observable emits
// something.
func (o Observable) Sample(notifier Observable) Observable {
	op := sampleOperator{notifier}
	return o.Lift(op.Call)
}
