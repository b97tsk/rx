package rx

import (
	"context"
)

type isEmptyOperator struct{}

func (op isEmptyOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
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

// IsEmpty creates an Observable that emits true if the source Observable
// emits no items, otherwise, it emits false.
func (o Observable) IsEmpty() Observable {
	op := isEmptyOperator{}
	return o.Lift(op.Call)
}
