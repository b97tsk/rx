package rx

import (
	"context"
)

type reduceOperator struct {
	source      Operator
	accumulator func(interface{}, interface{}, int) interface{}
	seed        interface{}
	hasSeed     bool
}

func (op reduceOperator) ApplyOptions(options []Option) Operator {
	for _, opt := range options {
		switch t := opt.(type) {
		case seedOption:
			op.seed = t.Value
			op.hasSeed = true
		default:
			panic(ErrUnsupportedOption)
		}
	}
	return op
}

func (op reduceOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	seed := op.seed
	hasSeed := op.hasSeed
	outerIndex := -1
	return op.source.Call(ctx, func(t Notification) {
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
	op := reduceOperator{
		source:      o.Op,
		accumulator: accumulator,
	}
	return Observable{op}
}
