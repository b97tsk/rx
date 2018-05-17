package rx

import (
	"context"
)

type reduceOperator struct {
	accumulator func(interface{}, interface{}, int) interface{}
	seed        interface{}
	hasSeed     bool
}

func (op reduceOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	var (
		seed       = op.seed
		hasSeed    = op.hasSeed
		outerIndex = -1
	)
	return source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			outerIndex++

			if hasSeed {
				seed = op.accumulator(seed, t.Value, outerIndex)
			} else {
				seed = t.Value
				hasSeed = true
			}

		case t.HasError:
			ob.Error(t.Value.(error))

		default:
			if hasSeed {
				ob.Next(seed)
			}
			ob.Complete()
		}
	})
}

// Reduce creates an Observable that applies an accumulator function over the
// source Observable, and emits the accumulated result when the source
// completes, given an optional seed value.
//
// Reduce combines together all values emitted on the source, using an
// accumulator function that knows how to join a new source value into the
// accumulation from the past.
func (o Observable) Reduce(accumulator func(interface{}, interface{}, int) interface{}) Observable {
	op := reduceOperator{accumulator: accumulator}
	return o.Lift(op.Call)
}
