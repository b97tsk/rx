package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// TakeWhile emits values emitted by the source so long as each value satisfies
// the given predicate, and then completes as soon as this predicate is not
// satisfied.
func TakeWhile(predicate func(interface{}, int) bool) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return takeWhileObservable{source, predicate}.Subscribe
	}
}

type takeWhileObservable struct {
	Source    rx.Observable
	Predicate func(interface{}, int) bool
}

func (obs takeWhileObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel)

	var observer rx.Observer

	sourceIndex := -1

	observer = func(t rx.Notification) {
		switch {
		case t.HasValue:
			sourceIndex++

			if obs.Predicate(t.Value, sourceIndex) {
				sink(t)
				break
			}

			observer = rx.Noop

			sink.Complete()

		default:
			sink(t)
		}
	}

	obs.Source.Subscribe(ctx, observer.Sink)
}
