package rx

import (
	"context"
)

type countOperator struct{}

func (op countOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	var count int
	return source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			count++
		case t.HasError:
			sink(t)
		default:
			sink.Next(count)
			sink.Complete()
		}
	})
}

// Count creates an Observable that counts the number of emissions on the
// source and emits that number when the source completes.
func (Operators) Count() OperatorFunc {
	return func(source Observable) Observable {
		op := countOperator{}
		return source.Lift(op.Call)
	}
}
