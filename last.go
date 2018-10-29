package rx

import (
	"context"
)

type lastOperator struct{}

func (op lastOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
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
func (Operators) Last() OperatorFunc {
	return func(source Observable) Observable {
		op := lastOperator{}
		return source.Lift(op.Call)
	}
}
