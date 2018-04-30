package rx

import (
	"context"
)

type excludeOperator struct {
	source    Operator
	predicate func(interface{}, int) bool
}

func (op excludeOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	outerIndex := -1
	return op.source.Call(ctx, ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			outerIndex++

			if !op.predicate(t.Value, outerIndex) {
				ob.Next(t.Value)
			}

		case t.HasError:
			ob.Error(t.Value.(error))

		default:
			ob.Complete()
		}
	}))
}

// Exclude creates an Observable that filter items emitted by the source
// Observable by only emitting those that do not satisfy a specified predicate.
func (o Observable) Exclude(predicate func(interface{}, int) bool) Observable {
	op := excludeOperator{
		source:    o.Op,
		predicate: predicate,
	}
	return Observable{op}
}
