package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// Skip creates an Observable that skips the first count items emitted by the
// source Observable.
func Skip(count int) rx.Operator {
	if count <= 0 {
		return noop
	}
	return func(source rx.Observable) rx.Observable {
		return skipObservable{source, count}.Subscribe
	}
}

type skipObservable struct {
	Source rx.Observable
	Count  int
}

func (obs skipObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	var observer rx.Observer

	count := obs.Count

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

	obs.Source.Subscribe(ctx, observer.Sink)
}
