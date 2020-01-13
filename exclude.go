package rx

import (
	"context"
)

type excludeObservable struct {
	Source    Observable
	Predicate func(interface{}, int) bool
}

func (obs excludeObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	var outerIndex = -1
	return obs.Source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			outerIndex++

			if !obs.Predicate(t.Value, outerIndex) {
				sink(t)
			}

		default:
			sink(t)
		}
	})
}

// Exclude creates an Observable that filter items emitted by the source
// Observable by only emitting those that do not satisfy a specified predicate.
func (Operators) Exclude(predicate func(interface{}, int) bool) OperatorFunc {
	return func(source Observable) Observable {
		return excludeObservable{source, predicate}.Subscribe
	}
}
