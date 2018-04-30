package rx

import (
	"context"
)

type countOperator struct {
	source    Operator
	predicate func(interface{}, int) bool
}

func (op countOperator) ApplyOptions(options []Option) Operator {
	for _, opt := range options {
		switch t := opt.(type) {
		case predicateOption:
			op.predicate = t.Value
		default:
			panic(ErrUnsupportedOption)
		}
	}
	return op
}

func (op countOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	outerIndex := -1
	count := 0
	return op.source.Call(ctx, ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			outerIndex++

			if op.predicate(t.Value, outerIndex) {
				count++
			}

		case t.HasError:
			ob.Error(t.Value.(error))

		default:
			ob.Next(count)
			ob.Complete()
		}
	}))
}

// Count creates an Observable that counts the number of emissions on the
// source and emits that number when the source completes.
//
// This operator takes an optional predicate function as argument, in which
// case the output emission will represent the number of source values that
// matched true with the predicate.
func (o Observable) Count() Observable {
	op := countOperator{
		source:    o.Op,
		predicate: defaultPredicate,
	}
	return Observable{op}
}
