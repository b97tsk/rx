package rx

import (
	"context"
)

type defaultIfEmptyOperator struct {
	defaultValue interface{}
}

func (op defaultIfEmptyOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	var hasValue bool
	return source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			hasValue = true
			t.Observe(ob)
		case t.HasError:
			t.Observe(ob)
		default:
			if !hasValue {
				ob.Next(op.defaultValue)
			}
			t.Observe(ob)
		}
	})
}

// DefaultIfEmpty creates an Observable that emits a given value if the source
// Observable completes without emitting any next value, otherwise mirrors the
// source Observable.
//
// If the source Observable turns out to be empty, then this operator will emit
// a default value.
func (o Observable) DefaultIfEmpty(defaultValue interface{}) Observable {
	op := defaultIfEmptyOperator{defaultValue}
	return o.Lift(op.Call)
}
