package rx

import (
	"context"
)

type countOperator struct{}

func (op countOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	var count int
	return source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			count++
		case t.HasError:
			ob.Error(t.Value.(error))
		default:
			ob.Next(count)
			ob.Complete()
		}
	})
}

// Count creates an Observable that counts the number of emissions on the
// source and emits that number when the source completes.
func (o Observable) Count() Observable {
	op := countOperator{}
	return o.Lift(op.Call)
}
