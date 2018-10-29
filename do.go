package rx

import (
	"context"
)

type doOperator struct {
	Sink Observer
}

func (op doOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	return source.Subscribe(ctx, func(t Notification) {
		op.Sink(t)
		sink(t)
	})
}

// Do creates an Observable that mirrors the source Observable, but perform
// a side effect before each emission.
func (Operators) Do(sink Observer) OperatorFunc {
	return func(source Observable) Observable {
		op := doOperator{sink}
		return source.Lift(op.Call)
	}
}
