package rx

import (
	"context"
)

type scanOperator struct {
	accumulator func(interface{}, interface{}, int) interface{}
	seed        interface{}
	hasSeed     bool
}

func (op scanOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
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

			ob.Next(seed)

		default:
			t.Observe(ob)
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
	op := scanOperator{accumulator: accumulator}
	return o.Lift(op.Call)
}
