package rx

import (
	"context"
)

type findOperator struct {
	predicate func(interface{}, int) bool
}

func (op findOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	var (
		outerIndex      = -1
		mutableObserver Observer
	)

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

	source.Subscribe(ctx, func(t Notification) { t.Observe(mutableObserver) })

	return ctx, cancel
}

// Find creates an Observable that emits only the first value emitted by the
// source Observable that meets some condition.
func (o Observable) Find(predicate func(interface{}, int) bool) Observable {
	op := findOperator{predicate}
	return o.Lift(op.Call)
}
