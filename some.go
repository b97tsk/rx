package rx

import (
	"context"
)

type someOperator struct {
	predicate func(interface{}, int) bool
}

func (op someOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
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
				sink.Next(true)
				sink.Complete()
				cancel()
			}

		case t.HasError:
			sink(t)
			cancel()

		default:
			sink.Next(false)
			sink.Complete()
			cancel()
		}
	}

	source.Subscribe(ctx, observer.Notify)

	return ctx, cancel
}

// Some creates an Observable that emits whether or not any item of the source
// satisfies the condition specified.
//
// Some emits true or false, then completes.
func (o Observable) Some(predicate func(interface{}, int) bool) Observable {
	op := someOperator{predicate}
	return o.Lift(op.Call)
}
