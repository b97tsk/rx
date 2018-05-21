package rx

import (
	"context"
)

type filterOperator struct {
	Predicate func(interface{}, int) bool
}

func (op filterOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	var outerIndex = -1
	return source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			outerIndex++

			if op.Predicate(t.Value, outerIndex) {
				sink(t)
			}

		default:
			sink(t)
		}
	})
}

// Filter creates an Observable that filter items emitted by the source
// Observable by only emitting those that satisfy a specified predicate.
func (o Observable) Filter(predicate func(interface{}, int) bool) Observable {
	op := filterOperator{predicate}
	return o.Lift(op.Call)
}
