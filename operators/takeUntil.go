package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// TakeUntil creates an Observable that emits the values emitted by the source
// Observable until a notifier Observable emits a value.
//
// TakeUntil lets values pass until a second Observable, notifier, emits
// something. Then, it completes.
func TakeUntil(notifier rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return takeUntilObservable{source, notifier}.Subscribe
	}
}

type takeUntilObservable struct {
	Source   rx.Observable
	Notifier rx.Observable
}

func (obs takeUntilObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel).Mutex()

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
