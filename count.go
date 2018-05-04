package rx

import (
	"context"
)

type countOperator struct {
	source Operator
}

func (op countOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	count := 0
	return op.source.Call(ctx, ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			count++
		case t.HasError:
			ob.Error(t.Value.(error))
		default:
			ob.Next(count)
			ob.Complete()
		}
	}))
}

// Count creates an Observable that counts the number of emissions on the
// source and emits that number when the source completes.
func (o Observable) Count() Observable {
	op := countOperator{o.Op}
	return Observable{op}
}
