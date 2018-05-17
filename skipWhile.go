package rx

import (
	"context"
)

type skipWhileOperator struct {
	predicate func(interface{}, int) bool
}

func (op skipWhileOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	var (
		outerIndex      = -1
		mutableObserver Observer
	)

	mutableObserver = func(t Notification) {
		switch {
		case t.HasValue:
			outerIndex++

			if !op.predicate(t.Value, outerIndex) {
				mutableObserver = ob
				ob.Next(t.Value)
			}

		case t.HasError:
			ob.Error(t.Value.(error))

		default:
			ob.Complete()
		}
	}

	return source.Subscribe(ctx, func(t Notification) { t.Observe(mutableObserver) })
}

// SkipWhile creates an Observable that skips all items emitted by the source
// Observable as long as a specified condition holds true, but emits all
// further source items as soon as the condition becomes false.
func (o Observable) SkipWhile(predicate func(interface{}, int) bool) Observable {
	op := skipWhileOperator{predicate}
	return o.Lift(op.Call)
}
