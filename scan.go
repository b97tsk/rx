package rx

import (
	"context"
)

type scanOperator struct {
	Accumulator func(interface{}, interface{}, int) interface{}
	Seed        interface{}
	HasSeed     bool
}

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
func (o Observable) Scan(accumulator func(interface{}, interface{}, int) interface{}) Observable {
	op := scanOperator{Accumulator: accumulator}
	return o.Lift(op.Call)
}
