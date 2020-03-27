package rx

import (
	"context"
)

type takeObservable struct {
	Source Observable
	Count  int
}

func (obs takeObservable) Subscribe(ctx context.Context, sink Observer) {
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
		obs := takeObservable{source, count}
		return Create(obs.Subscribe)
	}
}
