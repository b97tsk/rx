package rx

import (
	"context"
)

type findIndexOperator struct {
	predicate func(interface{}, int) bool
}

func (op findIndexOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
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
				ob.Next(outerIndex)
				ob.Complete()
				cancel()
			}

		default:
			t.Observe(ob)
			cancel()
		}
	}

	source.Subscribe(ctx, func(t Notification) { t.Observe(mutableObserver) })

	return ctx, cancel
}

// FindIndex creates an Observable that emits only the index of the first value
// emitted by the source Observable that meets some condition.
func (o Observable) FindIndex(predicate func(interface{}, int) bool) Observable {
	op := findIndexOperator{predicate}
	return o.Lift(op.Call)
}
