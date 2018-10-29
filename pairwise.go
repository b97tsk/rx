package rx

import (
	"context"
)

type pairwiseOperator struct{}

func (op pairwiseOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	var (
		prev    interface{}
		hasPrev bool
	)
	return source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			if hasPrev {
				sink.Next([]interface{}{prev, t.Value})
				prev = t.Value
			} else {
				prev = t.Value
				hasPrev = true
			}
		default:
			sink(t)
		}
	})
}

// Pairwise creates an Observable that groups pairs of consecutive emissions
// together and emits them as a slice of two values.
func (Operators) Pairwise() OperatorFunc {
	return func(source Observable) Observable {
		op := pairwiseOperator{}
		return source.Lift(op.Call)
	}
}
