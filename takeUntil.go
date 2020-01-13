package rx

import (
	"context"
)

type takeUntilObservable struct {
	Source   Observable
	Notifier Observable
}

func (obs takeUntilObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Mutex(Finally(sink, cancel))

	obs.Notifier.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			sink.Complete()
		default:
			sink(t)
		}
	})

	if ctx.Err() != nil {
		return Done()
	}

	obs.Source.Subscribe(ctx, sink)

	return ctx, cancel
}

// TakeUntil creates an Observable that emits the values emitted by the source
// Observable until a notifier Observable emits a value.
//
// TakeUntil lets values pass until a second Observable, notifier, emits
// something. Then, it completes.
func (Operators) TakeUntil(notifier Observable) OperatorFunc {
	return func(source Observable) Observable {
		return takeUntilObservable{source, notifier}.Subscribe
	}
}
