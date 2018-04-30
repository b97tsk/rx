package rx

import (
	"context"
)

type everyOperator struct {
	source    Operator
	predicate func(interface{}, int) bool
}

func (op everyOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	outerIndex := -1

	mutable := MutableObserver{}

	mutable.Observer = ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			outerIndex++

			if !op.predicate(t.Value, outerIndex) {
				mutable.Observer = NopObserver
				ob.Next(false)
				ob.Complete()
				cancel()
			}

		case t.HasError:
			ob.Error(t.Value.(error))
			cancel()

		default:
			ob.Next(true)
			ob.Complete()
			cancel()
		}
	})

	op.source.Call(ctx, &mutable)

	return ctx, cancel
}

// Every creates an Observable that emits whether or not every item of the source
// satisfies the condition specified.
//
// Every emits true or false, then completes.
func (o Observable) Every(predicate func(interface{}, int) bool) Observable {
	op := everyOperator{
		source:    o.Op,
		predicate: predicate,
	}
	return Observable{op}
}
