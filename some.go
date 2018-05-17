package rx

import (
	"context"
)

type someOperator struct {
	source    Operator
	predicate func(interface{}, int) bool
}

func (op someOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	outerIndex := -1

	var mutableObserver Observer

	mutableObserver = func(t Notification) {
		switch {
		case t.HasValue:
			outerIndex++

			if op.predicate(t.Value, outerIndex) {
				mutableObserver = NopObserver
				ob.Next(true)
				ob.Complete()
				cancel()
			}

		case t.HasError:
			ob.Error(t.Value.(error))
			cancel()

		default:
			ob.Next(false)
			ob.Complete()
			cancel()
		}
	}

	op.source.Call(ctx, func(t Notification) { t.Observe(mutableObserver) })

	return ctx, cancel
}

// Some creates an Observable that emits whether or not any item of the source
// satisfies the condition specified.
//
// Some emits true or false, then completes.
func (o Observable) Some(predicate func(interface{}, int) bool) Observable {
	op := someOperator{
		source:    o.Op,
		predicate: predicate,
	}
	return Observable{op}
}
