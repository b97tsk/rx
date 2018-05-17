package rx

import (
	"context"
)

type everyOperator struct {
	predicate func(interface{}, int) bool
}

func (op everyOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	var (
		outerIndex      = -1
		mutableObserver Observer
	)

	mutableObserver = func(t Notification) {
		switch {
		case t.HasValue:
			outerIndex++

			if !op.predicate(t.Value, outerIndex) {
				mutableObserver = NopObserver
				ob.Next(false)
				ob.Complete()
				cancel()
			}

		case t.HasError:
			t.Observe(ob)
			cancel()

		default:
			ob.Next(true)
			ob.Complete()
			cancel()
		}
	}

	source.Subscribe(ctx, func(t Notification) { t.Observe(mutableObserver) })

	return ctx, cancel
}

// Every creates an Observable that emits whether or not every item of the source
// satisfies the condition specified.
//
// Every emits true or false, then completes.
func (o Observable) Every(predicate func(interface{}, int) bool) Observable {
	op := everyOperator{predicate}
	return o.Lift(op.Call)
}
