package rx

import (
	"context"
)

type pairwiseOperator struct {
	source Operator
}

func (op pairwiseOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	prev := interface{}(nil)
	hasPrev := false
	return op.source.Call(ctx, ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			if hasPrev {
				ob.Next([]interface{}{prev, t.Value})
				prev = t.Value
			} else {
				prev = t.Value
				hasPrev = true
			}
		case t.HasError:
			ob.Error(t.Value.(error))
		default:
			ob.Complete()
		}
	}))
}

// Pairwise creates an Observable that groups pairs of consecutive emissions
// together and emits them as a slice of two values.
func (o Observable) Pairwise() Observable {
	op := pairwiseOperator{o.Op}
	return Observable{op}
}
