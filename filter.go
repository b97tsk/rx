package rx

import (
	"context"
)

type filterObservable struct {
	Source    Observable
	Predicate func(interface{}, int) bool
}

func (obs filterObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	var outerIndex = -1
	return obs.Source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			outerIndex++

			if obs.Predicate(t.Value, outerIndex) {
				sink(t)
			}

		default:
			sink(t)
		}
	})
}

// Filter creates an Observable that filter items emitted by the source
// Observable by only emitting those that satisfy a specified predicate.
func (Operators) Filter(predicate func(interface{}, int) bool) Operator {
	return func(source Observable) Observable {
		return filterObservable{source, predicate}.Subscribe
	}
}
