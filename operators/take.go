package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// Take emits only the first count values emitted by the source.
func Take(count int) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		if count <= 0 {
			return rx.Empty()
		}

		return takeObservable{source, count}.Subscribe
	}
}

type takeObservable struct {
	Source rx.Observable
	Count  int
}

func (obs takeObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel)

	var observer rx.Observer

	count := obs.Count

	observer = func(t rx.Notification) {
		sink(t)

		if t.HasValue {
			if count--; count == 0 {
				observer = rx.Noop

				sink.Complete()
			}
		}
	}

	obs.Source.Subscribe(ctx, observer.Sink)
}
