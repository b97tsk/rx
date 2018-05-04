package rx

import (
	"context"
)

type lastOperator struct {
	source Operator
}

func (op lastOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	lastValue := interface{}(nil)
	hasLastValue := false
	return op.source.Call(ctx, ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			lastValue = t.Value
			hasLastValue = true
		case t.HasError:
			ob.Error(t.Value.(error))
		default:
			if hasLastValue {
				ob.Next(lastValue)
				ob.Complete()
			} else {
				ob.Error(ErrEmpty)
			}
		}
	}))
}

// Last creates an Observable that emits only the last item emitted by the
// source Observable.
func (o Observable) Last() Observable {
	op := lastOperator{o.Op}
	return Observable{op}
}
