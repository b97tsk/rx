package rx

import (
	"context"
)

type lastOperator struct{}

func (op lastOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	var (
		lastValue    interface{}
		hasLastValue bool
	)
	return source.Subscribe(ctx, func(t Notification) {
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
	})
}

// Last creates an Observable that emits only the last item emitted by the
// source Observable.
func (o Observable) Last() Observable {
	op := lastOperator{}
	return o.Lift(op.Call)
}
