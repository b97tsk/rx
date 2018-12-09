package rx

import (
	"context"
)

type reduceOperator struct {
	Accumulator func(interface{}, interface{}, int) interface{}
	Seed        interface{}
	HasSeed     bool
}

func (op reduceOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	var (
		seed       = op.Seed
		hasSeed    = op.HasSeed
		outerIndex = -1
	)
	return source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			outerIndex++

			if hasSeed {
				seed = op.Accumulator(seed, t.Value, outerIndex)
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
func (Operators) Reduce(accumulator func(interface{}, interface{}, int) interface{}) OperatorFunc {
	return func(source Observable) Observable {
		op := reduceOperator{Accumulator: accumulator}
		return source.Lift(op.Call)
	}
}

// Fold creates an Observable that applies an accumulator function over
// the source Observable, and emits the accumulated result when the source
// completes, given an initial value.
//
// It's like Reduce, but you could specify an initial value.
func (Operators) Fold(initialValue interface{}, accumulator func(interface{}, interface{}, int) interface{}) OperatorFunc {
	return func(source Observable) Observable {
		op := reduceOperator{
			Accumulator: accumulator,
			Seed:        initialValue,
			HasSeed:     true,
		}
		return source.Lift(op.Call)
	}
}
