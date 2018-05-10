package rx

import (
	"context"
)

type sampleOperator struct {
	source   Operator
	notifier Observable
}

func (op sampleOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	try := cancellableLocker{}
	latestValue := interface{}(nil)
	hasLatestValue := false

	op.notifier.Subscribe(ctx, ObserverFunc(func(t Notification) {
		if t.HasError {
			ob.Error(t.Value.(error))
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
	}))

	op.source.Call(ctx, ObserverFunc(func(t Notification) {
		if try.Lock() {
			switch {
			case t.HasValue:
				latestValue = t.Value
				hasLatestValue = true
				try.Unlock()
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

// Sample creates an Observable that emits the most recently emitted value from
// the source Observable whenever another Observable, the notifier, emits.
//
// It's like SampleTime, but samples whenever the notifier Observable emits
// something.
func (o Observable) Sample(notifier Observable) Observable {
	op := sampleOperator{
		source:   o.Op,
		notifier: notifier,
	}
	return Observable{op}
}
