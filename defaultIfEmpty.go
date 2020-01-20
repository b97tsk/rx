package rx

import (
	"context"
)

// DefaultIfEmpty creates an Observable that emits a given value if the source
// Observable completes without emitting any next value, otherwise mirrors the
// source Observable.
//
// If the source Observable turns out to be empty, then this operator will emit
// a default value.
func (Operators) DefaultIfEmpty(defaultValue interface{}) Operator {
	return func(source Observable) Observable {
		return func(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
			var hasValue bool
			return source.Subscribe(ctx, func(t Notification) {
				switch {
				case t.HasValue:
					hasValue = true
					sink(t)
				case t.HasError:
					sink(t)
				default:
					if !hasValue {
						sink.Next(defaultValue)
					}
					sink(t)
				}
			})
		}
	}
}
