package rx

import (
	"context"
)

type skipWhileObservable struct {
	Source    Observable
	Predicate func(interface{}, int) bool
}

func (obs skipWhileObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	var (
		outerIndex = -1
		observer   Observer
	)

	observer = func(t Notification) {
		switch {
		case t.HasValue:
			outerIndex++

			if !obs.Predicate(t.Value, outerIndex) {
				observer = sink
				sink(t)
			}

		default:
			sink(t)
		}
	}

	return obs.Source.Subscribe(ctx, observer.Notify)
}

// SkipWhile creates an Observable that skips all items emitted by the source
// Observable as long as a specified condition holds true, but emits all
// further source items as soon as the condition becomes false.
func (Operators) SkipWhile(predicate func(interface{}, int) bool) Operator {
	return func(source Observable) Observable {
		return skipWhileObservable{source, predicate}.Subscribe
	}
}
