package rx

import (
	"context"
)

type findOperator struct {
	source    Operator
	predicate func(interface{}, int) bool
}

func (op findOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	outerIndex := -1

	var mutableObserver Observer

	mutableObserver = func(t Notification) {
		switch {
		case t.HasValue:
			outerIndex++

			if op.predicate(t.Value, outerIndex) {
				mutableObserver = NopObserver
				ob.Next(t.Value)
				ob.Complete()
				cancel()
			}

		case t.HasError:
			ob.Error(t.Value.(error))
			cancel()

		default:
			ob.Complete()
			cancel()
		}
	}

	op.source.Call(ctx, func(t Notification) { t.Observe(mutableObserver) })

	return ctx, cancel
}

// Find creates an Observable that emits only the first value emitted by the
// source Observable that meets some condition.
func (o Observable) Find(predicate func(interface{}, int) bool) Observable {
	op := findOperator{
		source:    o.Op,
		predicate: predicate,
	}
	return Observable{op}
}
