package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

type skipObservable struct {
	Source rx.Observable
	Count  int
}

func (obs skipObservable) Subscribe(ctx context.Context, sink rx.Observer) (context.Context, context.CancelFunc) {
	var (
		count    = obs.Count
		observer rx.Observer
	)

	observer = func(t rx.Notification) {
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
func Skip(count int) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		if count <= 0 {
			return source
		}
		return skipObservable{source, count}.Subscribe
	}
}
