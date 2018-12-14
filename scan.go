package rx

import (
	"context"
)

// A ScanConfigure is a configure for Scan.
type ScanConfigure struct {
	Accumulator func(interface{}, interface{}, int) interface{}
	Seed        interface{}
	HasSeed     bool
}

// MakeFunc creates an OperatorFunc from this type.
func (conf ScanConfigure) MakeFunc() OperatorFunc {
	return MakeFunc(scanOperator(conf).Call)
}

type scanOperator ScanConfigure

func (op scanOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
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
		op := scanOperator{Accumulator: accumulator}
		return source.Lift(op.Call)
	}
}
