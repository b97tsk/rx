package rx

import (
	"context"
)

type defaultIfEmptyOperator struct {
	source       Operator
	defaultValue interface{}
}

func (op defaultIfEmptyOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	hasValue := false
	return op.source.Call(ctx, ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			hasValue = true
			ob.Next(t.Value)
		case t.HasError:
			ob.Error(t.Value.(error))
		default:
			if !hasValue {
				ob.Next(op.defaultValue)
			}
			ob.Complete()
		}
	}))
}

// DefaultIfEmpty creates an Observable that emits a given value if the source
// Observable completes without emitting any next value, otherwise mirrors the
// source Observable.
//
// If the source Observable turns out to be empty, then this operator will emit
// a default value.
func (o Observable) DefaultIfEmpty(defaultValue interface{}) Observable {
	op := defaultIfEmptyOperator{
		source:       o.Op,
		defaultValue: defaultValue,
	}
	return Observable{op}
}
