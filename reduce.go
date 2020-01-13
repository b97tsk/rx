package rx

import (
	"context"
)

type reduceObservable struct {
	Source      Observable
	Accumulator func(interface{}, interface{}, int) interface{}
	Seed        interface{}
	HasSeed     bool
}

func (obs reduceObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	var (
		seed       = obs.Seed
		hasSeed    = obs.HasSeed
		outerIndex = -1
	)
	return obs.Source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			outerIndex++

			if hasSeed {
				seed = obs.Accumulator(seed, t.Value, outerIndex)
			} else {
				seed = t.Value
				hasSeed = true
			}

		case t.HasError:
			sink(t)

		default:
			if hasSeed {
				sink.Next(seed)
			}
			sink(t)
		}
	})
}

// Reduce creates an Observable that applies an accumulator function over
// the source Observable, and emits the accumulated result when the source
// completes.
//
// It's like Fold, but no need to specify an initial value.
func (Operators) Reduce(accumulator func(interface{}, interface{}, int) interface{}) Operator {
	return func(source Observable) Observable {
		return reduceObservable{
			Source:      source,
			Accumulator: accumulator,
		}.Subscribe
	}
}

// Fold creates an Observable that applies an accumulator function over
// the source Observable, and emits the accumulated result when the source
// completes, given an initial value.
//
// It's like Reduce, but you could specify an initial value.
func (Operators) Fold(initialValue interface{}, accumulator func(interface{}, interface{}, int) interface{}) Operator {
	return func(source Observable) Observable {
		return reduceObservable{
			Source:      source,
			Accumulator: accumulator,
			Seed:        initialValue,
			HasSeed:     true,
		}.Subscribe
	}
}
