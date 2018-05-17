package rx

import (
	"context"
)

type filterOperator struct {
	source    Operator
	predicate func(interface{}, int) bool
}

func (op filterOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	outerIndex := -1
	return op.source.Call(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			outerIndex++

			if op.predicate(t.Value, outerIndex) {
				ob.Next(t.Value)
			}

		case t.HasError:
			ob.Error(t.Value.(error))

		default:
			ob.Complete()
		}
	})
}

// Filter creates an Observable that filter items emitted by the source
// Observable by only emitting those that satisfy a specified predicate.
func (o Observable) Filter(predicate func(interface{}, int) bool) Observable {
	op := filterOperator{
		source:    o.Op,
		predicate: predicate,
	}
	return Observable{op}
}
