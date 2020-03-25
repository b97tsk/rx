package rx

import (
	"context"
)

type firstObservable struct {
	Source Observable
}

func (obs firstObservable) Subscribe(parent context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx := NewContext(parent)

	sink = DoAtLast(sink, ctx.AtLast)

	var observer Observer
	observer = func(t Notification) {
		switch {
		case t.HasValue:
			observer = NopObserver
			sink(t)
			sink.Complete()
		case t.HasError:
			sink(t)
		default:
			sink.Error(ErrEmpty)
		}
	}
	obs.Source.Subscribe(ctx, observer.Notify)

	return ctx, ctx.Cancel
}

// First creates an Observable that emits only the first value (or the first
// value that meets some condition) emitted by the source Observable.
func (Operators) First() Operator {
	return func(source Observable) Observable {
		return firstObservable{source}.Subscribe
	}
}
