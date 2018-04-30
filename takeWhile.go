package rx

import (
	"context"
)

type takeWhileOperator struct {
	source    Operator
	predicate func(interface{}, int) bool
}

func (op takeWhileOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	outerIndex := -1

	mutable := MutableObserver{}

	mutable.Observer = ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			outerIndex++

			if op.predicate(t.Value, outerIndex) {
				ob.Next(t.Value)
				break
			}

			mutable.Observer = NopObserver
			ob.Complete()
			cancel()

		case t.HasError:
			ob.Error(t.Value.(error))
			cancel()

		default:
			ob.Complete()
			cancel()
		}
	})

	op.source.Call(ctx, &mutable)

	return ctx, cancel
}

// TakeWhile creates an Observable that emits values emitted by the source
// Observable so long as each value satisfies the given predicate, and then
// completes as soon as this predicate is not satisfied.
//
// TakeWhile takes values from the source only while they pass the condition
// given. When the first value does not satisfy, it completes.
func (o Observable) TakeWhile(predicate func(interface{}, int) bool) Observable {
	op := takeWhileOperator{
		source:    o.Op,
		predicate: predicate,
	}
	return Observable{op}
}
