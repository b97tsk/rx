package rx

import (
	"context"
)

type isEmptyOperator struct {
	source Operator
}

func (op isEmptyOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	mutable := MutableObserver{}

	mutable.Observer = ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			mutable.Observer = NopObserver
			ob.Next(false)
			ob.Complete()
			cancel()
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

// IsEmpty creates an Observable that emits true if the source Observable
// emits no items, otherwise, it emits false.
func (o Observable) IsEmpty() Observable {
	op := isEmptyOperator{o.Op}
	return Observable{op}
}
