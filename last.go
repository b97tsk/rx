package rx

import (
	"context"
)

type lastObservable struct {
	Source Observable
}

func (obs lastObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	var (
		lastValue    interface{}
		hasLastValue bool
	)
	return obs.Source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			lastValue = t.Value
			hasLastValue = true
		case t.HasError:
			sink(t)
		default:
			if hasLastValue {
				sink.Next(lastValue)
				sink.Complete()
			} else {
				sink.Error(ErrEmpty)
			}
		}
	})
}

// Last creates an Observable that emits only the last item emitted by the
// source Observable.
func (Operators) Last() Operator {
	return func(source Observable) Observable {
		return lastObservable{source}.Subscribe
	}
}
