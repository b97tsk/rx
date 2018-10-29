package rx

import (
	"context"
)

type materializeOperator struct{}

func (op materializeOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	return source.Subscribe(ctx, func(t Notification) {
		sink.Next(t)

		if t.HasValue {
			return
		}

		sink.Complete()
	})
}

// Materialize creates an Observable that represents all of the notifications
// from the source Observable as Next emissions, and then completes.
//
// Materialize wraps Next, Error and Complete emissions in Notification objects,
// emitted as Next on the output Observable.
func (Operators) Materialize() OperatorFunc {
	return func(source Observable) Observable {
		op := materializeOperator{}
		return source.Lift(op.Call)
	}
}
