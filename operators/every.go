package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// Every creates an Observable that emits whether or not every item of the source
// satisfies a specified predicate.
//
// Every emits true or false, then completes.
func Every(predicate func(interface{}, int) bool) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return everyObservable{source, predicate}.Subscribe
	}
}

type everyObservable struct {
	Source    rx.Observable
	Predicate func(interface{}, int) bool
}

func (obs everyObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel)

	var observer rx.Observer

	sourceIndex := -1

	observer = func(t rx.Notification) {
		switch {
		case t.HasValue:
			sourceIndex++

			if !obs.Predicate(t.Value, sourceIndex) {
				observer = rx.Noop

				sink.Next(false)
				sink.Complete()
			}

		case t.HasError:
			sink(t)

		default:
			sink.Next(true)
			sink.Complete()
		}
	}

	obs.Source.Subscribe(ctx, observer.Sink)
}
