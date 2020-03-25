package rx

import (
	"context"
)

type takeUntilObservable struct {
	Source   Observable
	Notifier Observable
}

func (obs takeUntilObservable) Subscribe(parent context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx := NewContext(parent)

	sink = Mutex(DoAtLast(sink, ctx.AtLast))

	obs.Notifier.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			sink.Complete()
		default:
			sink(t)
		}
	})

	if ctx.Err() != nil {
		return ctx, ctx.Cancel
	}

	obs.Source.Subscribe(ctx, sink)

	return ctx, ctx.Cancel
}

// TakeUntil creates an Observable that emits the values emitted by the source
// Observable until a notifier Observable emits a value.
//
// TakeUntil lets values pass until a second Observable, notifier, emits
// something. Then, it completes.
func (Operators) TakeUntil(notifier Observable) Operator {
	return func(source Observable) Observable {
		return takeUntilObservable{source, notifier}.Subscribe
	}
}
