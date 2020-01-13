package rx

import (
	"context"
)

type defaultIfEmptyObservable struct {
	Source  Observable
	Default interface{}
}

func (obs defaultIfEmptyObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	var hasValue bool
	return obs.Source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			hasValue = true
			sink(t)
		case t.HasError:
			sink(t)
		default:
			if !hasValue {
				sink.Next(obs.Default)
			}
			sink(t)
		}
	})
}

// DefaultIfEmpty creates an Observable that emits a given value if the source
// Observable completes without emitting any next value, otherwise mirrors the
// source Observable.
//
// If the source Observable turns out to be empty, then this operator will emit
// a default value.
func (Operators) DefaultIfEmpty(defaultValue interface{}) OperatorFunc {
	return func(source Observable) Observable {
		return defaultIfEmptyObservable{source, defaultValue}.Subscribe
	}
}
