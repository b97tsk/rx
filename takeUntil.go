package rx

import (
	"context"
)

type takeUntilObservable struct {
	Source   Observable
	Notifier Observable
}

func (obs takeUntilObservable) Subscribe(ctx context.Context, sink Observer) {
	sink = Mutex(sink)

	obs.Notifier.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			sink.Complete()
		default:
			sink(t)
		}
	})

	if ctx.Err() != nil {
		return
	}

	obs.Source.Subscribe(ctx, sink)
}

// TakeUntil creates an Observable that emits the values emitted by the source
// Observable until a notifier Observable emits a value.
//
// TakeUntil lets values pass until a second Observable, notifier, emits
// something. Then, it completes.
func (Operators) TakeUntil(notifier Observable) Operator {
	return func(source Observable) Observable {
		obs := takeUntilObservable{source, notifier}
		return Create(obs.Subscribe)
	}
}
