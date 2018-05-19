package rx

import (
	"context"
)

type findIndexOperator struct {
	predicate func(interface{}, int) bool
}

func (op findIndexOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	var (
		outerIndex = -1
		observer   Observer
	)

	observer = func(t Notification) {
		switch {
		case t.HasValue:
			outerIndex++

			if op.predicate(t.Value, outerIndex) {
				observer = NopObserver
				sink.Next(outerIndex)
				sink.Complete()
				cancel()
			}

		default:
			sink(t)
			cancel()
		}
	}

	source.Subscribe(ctx, observer.Notify)

	return ctx, cancel
}

// FindIndex creates an Observable that emits only the index of the first value
// emitted by the source Observable that meets some condition.
func (o Observable) FindIndex(predicate func(interface{}, int) bool) Observable {
	op := findIndexOperator{predicate}
	return o.Lift(op.Call)
}
