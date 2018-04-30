package rx

import (
	"context"
)

type findOperator struct {
	source    Operator
	predicate func(interface{}, int) bool
}

func (op findOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	outerIndex := -1

	mutable := MutableObserver{}

	mutable.Observer = ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			outerIndex++

			if op.predicate(t.Value, outerIndex) {
				mutable.Observer = NopObserver
				ob.Next(t.Value)
				ob.Complete()
				cancel()
			}

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

// First creates an Observable that emits only the first value (or the first
// value that meets some condition) emitted by the source Observable.
func (o Observable) First() Observable {
	op := findOperator{
		source:    o.Op,
		predicate: defaultPredicate,
	}
	return Observable{op}
}

// Find creates an Observable that emits only the first value emitted by the
// source Observable that meets some condition.
func (o Observable) Find(predicate func(interface{}, int) bool) Observable {
	op := findOperator{
		source:    o.Op,
		predicate: predicate,
	}
	return Observable{op}
}
