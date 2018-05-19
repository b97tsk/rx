package rx

import (
	"context"
)

type isEmptyOperator struct{}

func (op isEmptyOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	var observer Observer

	observer = func(t Notification) {
		switch {
		case t.HasValue:
			observer = NopObserver
			sink.Next(false)
			sink.Complete()
			cancel()
		case t.HasError:
			sink(t)
			cancel()
		default:
			sink.Next(true)
			sink.Complete()
			cancel()
		}
	}

	source.Subscribe(ctx, observer.Notify)

	return ctx, cancel
}

// IsEmpty creates an Observable that emits true if the source Observable
// emits no items, otherwise, it emits false.
func (o Observable) IsEmpty() Observable {
	op := isEmptyOperator{}
	return o.Lift(op.Call)
}
