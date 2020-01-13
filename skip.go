package rx

import (
	"context"
)

type skipObservable struct {
	Source Observable
	Count  int
}

func (obs skipObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	var (
		count    = obs.Count
		observer Observer
	)

	observer = func(t Notification) {
		switch {
		case t.HasValue:
			if count > 1 {
				count--
			} else {
				observer = sink
			}
		default:
			sink(t)
		}
	}

	return obs.Source.Subscribe(ctx, observer.Notify)
}

// Skip creates an Observable that skips the first count items emitted by the
// source Observable.
func (Operators) Skip(count int) Operator {
	return func(source Observable) Observable {
		if count <= 0 {
			return source
		}
		return skipObservable{source, count}.Subscribe
	}
}
