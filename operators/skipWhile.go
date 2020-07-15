package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

type skipWhileObservable struct {
	Source    rx.Observable
	Predicate func(interface{}, int) bool
}

func (obs skipWhileObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	var observer rx.Observer

	sourceIndex := -1

	observer = func(t rx.Notification) {
		switch {
		case t.HasValue:
			sourceIndex++

			if !obs.Predicate(t.Value, sourceIndex) {
				observer = sink
				sink(t)
			}

		default:
			sink(t)
		}
	}

	obs.Source.Subscribe(ctx, observer.Sink)
}

// SkipWhile creates an Observable that skips all items emitted by the source
// Observable as long as a specified condition holds true, but emits all
// further source items as soon as the condition becomes false.
func SkipWhile(predicate func(interface{}, int) bool) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return skipWhileObservable{source, predicate}.Subscribe
	}
}
