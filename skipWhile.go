package rx

import (
	"context"
)

type skipWhileOperator struct {
	source    Operator
	predicate func(interface{}, int) bool
}

func (op skipWhileOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	outerIndex := -1

	mutable := MutableObserver{}

	mutable.Observer = ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			outerIndex++

			if !op.predicate(t.Value, outerIndex) {
				mutable.Observer = ob
				ob.Next(t.Value)
			}

		case t.HasError:
			ob.Error(t.Value.(error))

		default:
			ob.Complete()
		}
	})

	return op.source.Call(ctx, &mutable)
}

// SkipWhile creates an Observable that skips all items emitted by the source
// Observable as long as a specified condition holds true, but emits all
// further source items as soon as the condition becomes false.
func (o Observable) SkipWhile(predicate func(interface{}, int) bool) Observable {
	op := skipWhileOperator{
		source:    o.Op,
		predicate: predicate,
	}
	return Observable{op}
}
