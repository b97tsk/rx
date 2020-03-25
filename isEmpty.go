package rx

import (
	"context"
)

type isEmptyObservable struct {
	Source Observable
}

func (obs isEmptyObservable) Subscribe(parent context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx := NewContext(parent)

	sink = DoAtLast(sink, ctx.AtLast)

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

	return ctx, ctx.Cancel
}

// IsEmpty creates an Observable that emits true if the source Observable
// emits no items, otherwise, it emits false.
func (Operators) IsEmpty() Operator {
	return func(source Observable) Observable {
		return isEmptyObservable{source}.Subscribe
	}
}
