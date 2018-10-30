package rx

import (
	"context"
)

// ScanOperator is an operator type.
type ScanOperator struct {
	Accumulator func(interface{}, interface{}, int) interface{}
	Seed        interface{}
	HasSeed     bool
}

// MakeFunc creates an OperatorFunc from this operator.
func (op ScanOperator) MakeFunc() OperatorFunc {
	return MakeFunc(op.Call)
}

// Call invokes an execution of this operator.
func (op ScanOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
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

			sink.Next(seed)

		default:
			sink(t)
		}
	})
}

// Scan creates an Observable that applies an accumulator function over the
// source Observable, and emits each intermediate result, with an optional
// seed value.
//
// It's like Reduce, but emits the current accumulation whenever the source
// emits a value.
func (Operators) Scan(accumulator func(interface{}, interface{}, int) interface{}) OperatorFunc {
	return func(source Observable) Observable {
		op := ScanOperator{Accumulator: accumulator}
		return source.Lift(op.Call)
	}
}
