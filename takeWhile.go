package rx

import (
	"context"
)

type takeWhileOperator struct {
	predicate func(interface{}, int) bool
}

func (op takeWhileOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
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
				t.Observe(ob)
				break
			}

			mutableObserver = NopObserver
			ob.Complete()
			cancel()

		default:
			t.Observe(ob)
			cancel()
		}
	}

	source.Subscribe(ctx, func(t Notification) { t.Observe(mutableObserver) })

	return ctx, cancel
}

// TakeWhile creates an Observable that emits values emitted by the source
// Observable so long as each value satisfies the given predicate, and then
// completes as soon as this predicate is not satisfied.
//
// TakeWhile takes values from the source only while they pass the condition
// given. When the first value does not satisfy, it completes.
func (o Observable) TakeWhile(predicate func(interface{}, int) bool) Observable {
	op := takeWhileOperator{predicate}
	return o.Lift(op.Call)
}
