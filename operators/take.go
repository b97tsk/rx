package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

type takeObservable struct {
	Source rx.Observable
	Count  int
}

func (obs takeObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	var (
		count    = obs.Count
		observer rx.Observer
	)

	observer = func(t rx.Notification) {
		switch {
		case t.HasValue:
			if count > 1 {
				count--
				sink(t)
			} else {
				observer = rx.Noop
				sink(t)
				sink.Complete()
			}
		default:
			sink(t)
		}
	}

	obs.Source.Subscribe(ctx, observer.Sink)
}

// Take creates an Observable that emits only the first count values emitted
// by the source Observable.
//
// Take takes the first count values from the source, then completes.
func Take(count int) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		if count <= 0 {
			return rx.Empty()
		}
		obs := takeObservable{source, count}
		return rx.Create(obs.Subscribe)
	}
}
