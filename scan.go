package rx

import (
	"context"
)

type scanOperator struct {
	source      Operator
	accumulator func(interface{}, interface{}, int) interface{}
	seed        interface{}
	hasSeed     bool
}

func (op scanOperator) ApplyOptions(options []Option) Operator {
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

func (op scanOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	seed := op.seed
	hasSeed := op.hasSeed
	outerIndex := -1
	return op.source.Call(ctx, ObserverFunc(func(t Notification) {
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

		case t.HasError:
			ob.Error(t.Value.(error))

		default:
			ob.Complete()
		}
	}))
}

// Scan creates an Observable that applies an accumulator function over the
// source Observable, and emits each intermediate result, with an optional
// seed value.
//
// It's like Reduce, but emits the current accumulation whenever the source
// emits a value.
func (o Observable) Scan(accumulator func(interface{}, interface{}, int) interface{}) Observable {
	op := scanOperator{
		source:      o.Op,
		accumulator: accumulator,
	}
	return Observable{op}
}
