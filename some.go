package rx

import (
	"context"
)

type someOperator struct {
	Predicate func(interface{}, int) bool
}

func (op someOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
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
				sink.Next(true)
				sink.Complete()
			}

		case t.HasError:
			sink(t)

		default:
			sink.Next(false)
			sink.Complete()
		}
	}

	source.Subscribe(ctx, observer.Notify)

	return ctx, cancel
}

// Some creates an Observable that emits whether or not any item of the source
// satisfies the condition specified.
//
// Some emits true or false, then completes.
func (Operators) Some(predicate func(interface{}, int) bool) OperatorFunc {
	return func(source Observable) Observable {
		op := someOperator{predicate}
		return source.Lift(op.Call)
	}
}
