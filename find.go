package rx

import (
	"context"
)

type findOperator struct {
	Predicate func(interface{}, int) bool
}

func (op findOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

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
				sink(t)
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

// Find creates an Observable that emits only the first value emitted by the
// source Observable that meets some condition.
func (Operators) Find(predicate func(interface{}, int) bool) OperatorFunc {
	return func(source Observable) Observable {
		op := findOperator{predicate}
		return source.Lift(op.Call)
	}
}
