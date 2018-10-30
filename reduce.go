package rx

import (
	"context"
)

// ReduceOperator is an operator type.
type ReduceOperator struct {
	Accumulator func(interface{}, interface{}, int) interface{}
	Seed        interface{}
	HasSeed     bool
}

// MakeFunc creates an OperatorFunc from this operator.
func (op ReduceOperator) MakeFunc() OperatorFunc {
	return MakeFunc(op.Call)
}

// Call invokes an execution of this operator.
func (op ReduceOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
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

// Reduce creates an Observable that applies an accumulator function over the
// source Observable, and emits the accumulated result when the source
// completes, given an optional seed value.
//
// Reduce combines together all values emitted on the source, using an
// accumulator function that knows how to join a new source value into the
// accumulation from the past.
func (Operators) Reduce(accumulator func(interface{}, interface{}, int) interface{}) OperatorFunc {
	return func(source Observable) Observable {
		op := ReduceOperator{Accumulator: accumulator}
		return source.Lift(op.Call)
	}
}
