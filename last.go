package rx

import (
	"context"
)

type lastOperator struct {
	source    Operator
	predicate func(interface{}, int) bool
}

func (op lastOperator) ApplyOptions(options []Option) Operator {
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

func (op lastOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	lastValue := interface{}(nil)
	hasLastValue := false
	outerIndex := -1
	return op.source.Call(ctx, ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			outerIndex++

			if op.predicate(t.Value, outerIndex) {
				lastValue = t.Value
				hasLastValue = true
			}

		case t.HasError:
			ob.Error(t.Value.(error))

		default:
			if hasLastValue {
				ob.Next(lastValue)
			}
			ob.Complete()
		}
	}))
}

// Last creates an Observable that emits only the last item emitted by the
// source Observable.
//
// It optionally takes a predicate function as a parameter, in which case,
// rather than emitting the last item from the source Observable, the resulting
// Observable will emit the last item from the source Observable that satisfies
// the predicate.
func (o Observable) Last() Observable {
	op := lastOperator{
		source:    o.Op,
		predicate: defaultPredicate,
	}
	return Observable{op}
}
