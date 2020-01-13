package rx

import (
	"context"
)

type isEmptyObservable struct {
	Source Observable
}

func (obs isEmptyObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	var observer Observer
	observer = func(t Notification) {
		switch {
		case t.HasValue:
			observer = NopObserver
			sink.Next(false)
			sink.Complete()
		case t.HasError:
			sink(t)
		default:
			sink.Next(true)
			sink.Complete()
		}
	}
	obs.Source.Subscribe(ctx, observer.Notify)

	return ctx, cancel
}

// IsEmpty creates an Observable that emits true if the source Observable
// emits no items, otherwise, it emits false.
func (Operators) IsEmpty() OperatorFunc {
	return func(source Observable) Observable {
		return isEmptyObservable{source}.Subscribe
	}
}
