package rx

import (
	"context"
)

type isEmptyOperator struct {
	source Operator
}

func (op isEmptyOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	var mutableObserver Observer

	mutableObserver = func(t Notification) {
		switch {
		case t.HasValue:
			mutableObserver = NopObserver
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
	}

	op.source.Call(ctx, func(t Notification) { t.Observe(mutableObserver) })

	return ctx, cancel
}

// IsEmpty creates an Observable that emits true if the source Observable
// emits no items, otherwise, it emits false.
func (o Observable) IsEmpty() Observable {
	op := isEmptyOperator{o.Op}
	return Observable{op}
}
