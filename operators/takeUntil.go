package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

type takeUntilObservable struct {
	Source   rx.Observable
	Notifier rx.Observable
}

func (obs takeUntilObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	sink = sink.Mutex()

	obs.Notifier.Subscribe(ctx, func(t rx.Notification) {
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
func TakeUntil(notifier rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := takeUntilObservable{source, notifier}
		return rx.Create(obs.Subscribe)
	}
}
