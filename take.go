package rx

import (
	"context"
)

type takeObservable struct {
	Source Observable
	Count  int
}

func (obs takeObservable) Subscribe(parent context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx := NewContext(parent)

	sink = DoAtLast(sink, ctx.AtLast)

	var (
		count    = obs.Count
		observer Observer
	)

	observer = func(t Notification) {
		switch {
		case t.HasValue:
			if count > 1 {
				count--
				sink(t)
			} else {
				observer = NopObserver
				sink(t)
				sink.Complete()
			}
		default:
			sink(t)
		}
	}

	obs.Source.Subscribe(ctx, observer.Notify)

	return ctx, ctx.Cancel
}

// Take creates an Observable that emits only the first count values emitted
// by the source Observable.
//
// Take takes the first count values from the source, then completes.
func (Operators) Take(count int) Operator {
	return func(source Observable) Observable {
		if count <= 0 {
			return Empty()
		}
		return takeObservable{source, count}.Subscribe
	}
}
