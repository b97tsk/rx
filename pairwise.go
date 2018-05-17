package rx

import (
	"context"
)

type pairwiseOperator struct{}

func (op pairwiseOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	var (
		prev    interface{}
		hasPrev bool
	)
	return source.Subscribe(ctx, func(t Notification) {
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
	})
}

// Pairwise creates an Observable that groups pairs of consecutive emissions
// together and emits them as a slice of two values.
func (o Observable) Pairwise() Observable {
	op := pairwiseOperator{}
	return o.Lift(op.Call)
}
