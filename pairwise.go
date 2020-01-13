package rx

import (
	"context"
)

type pairwiseObservable struct {
	Source Observable
}

func (obs pairwiseObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	var (
		prev    interface{}
		hasPrev bool
	)
	return obs.Source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			if hasPrev {
				sink.Next([]interface{}{prev, t.Value})
				prev = t.Value
			} else {
				prev = t.Value
				hasPrev = true
			}
		default:
			sink(t)
		}
	})
}

// Pairwise creates an Observable that groups pairs of consecutive emissions
// together and emits them as a slice of two values.
func (Operators) Pairwise() Operator {
	return func(source Observable) Observable {
		return pairwiseObservable{source}.Subscribe
	}
}
