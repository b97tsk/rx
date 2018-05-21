package rx

import (
	"context"
)

type skipWhileOperator struct {
	Predicate func(interface{}, int) bool
}

func (op skipWhileOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	var (
		outerIndex = -1
		observer   Observer
	)

	observer = func(t Notification) {
		switch {
		case t.HasValue:
			outerIndex++

			if !op.Predicate(t.Value, outerIndex) {
				observer = sink
				sink(t)
			}

		default:
			sink(t)
		}
	}

	return source.Subscribe(ctx, observer.Notify)
}

// SkipWhile creates an Observable that skips all items emitted by the source
// Observable as long as a specified condition holds true, but emits all
// further source items as soon as the condition becomes false.
func (o Observable) SkipWhile(predicate func(interface{}, int) bool) Observable {
	op := skipWhileOperator{predicate}
	return o.Lift(op.Call)
}
