package rx

import (
	"context"
)

type excludeOperator struct {
	Predicate func(interface{}, int) bool
}

func (op excludeOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	var outerIndex = -1
	return source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			outerIndex++

			if !op.Predicate(t.Value, outerIndex) {
				sink(t)
			}

		default:
			sink(t)
		}
	})
}

// Exclude creates an Observable that filter items emitted by the source
// Observable by only emitting those that do not satisfy a specified predicate.
func (Operators) Exclude(predicate func(interface{}, int) bool) OperatorFunc {
	return func(source Observable) Observable {
		op := excludeOperator{predicate}
		return source.Lift(op.Call)
	}
}
