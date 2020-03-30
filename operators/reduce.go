package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

type reduceObservable struct {
	Source      rx.Observable
	Accumulator func(interface{}, interface{}, int) interface{}
	Seed        interface{}
	HasSeed     bool
}

func (obs reduceObservable) Subscribe(ctx context.Context, sink rx.Observer) (context.Context, context.CancelFunc) {
	var (
		seed        = obs.Seed
		hasSeed     = obs.HasSeed
		sourceIndex = -1
	)
	return obs.Source.Subscribe(ctx, func(t rx.Notification) {
		switch {
		case t.HasValue:
			sourceIndex++

			if hasSeed {
				seed = obs.Accumulator(seed, t.Value, sourceIndex)
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
func Reduce(accumulator func(interface{}, interface{}, int) interface{}) rx.Operator {
	return func(source rx.Observable) rx.Observable {
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
func Fold(initialValue interface{}, accumulator func(interface{}, interface{}, int) interface{}) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return reduceObservable{
			Source:      source,
			Accumulator: accumulator,
			Seed:        initialValue,
			HasSeed:     true,
		}.Subscribe
	}
}
