package rx

import (
	"context"
)

type findIndexOperator struct {
	Predicate func(interface{}, int) bool
}

func (op findIndexOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	var (
		outerIndex = -1
		observer   Observer
	)

	observer = func(t Notification) {
		switch {
		case t.HasValue:
			outerIndex++

			if op.Predicate(t.Value, outerIndex) {
				observer = NopObserver
				sink.Next(outerIndex)
				sink.Complete()
			}

		default:
			sink(t)
		}
	}

	source.Subscribe(ctx, observer.Notify)

	return ctx, cancel
}

// FindIndex creates an Observable that emits only the index of the first value
// emitted by the source Observable that meets some condition.
func (Operators) FindIndex(predicate func(interface{}, int) bool) OperatorFunc {
	return func(source Observable) Observable {
		op := findIndexOperator{predicate}
		return source.Lift(op.Call)
	}
}
