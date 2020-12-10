package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// DefaultIfEmpty creates an Observable that emits a given value if the source
// Observable completes without emitting any next value, otherwise mirrors the
// source Observable.
//
// If the source Observable turns out to be empty, then this operator will emit
// a default value.
func DefaultIfEmpty(def interface{}) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return func(ctx context.Context, sink rx.Observer) {
			var hasValue bool

			source.Subscribe(ctx, func(t rx.Notification) {
				switch {
				case t.HasValue:
					hasValue = true

					sink(t)

				case t.HasError:
					sink(t)

				default:
					if !hasValue {
						sink.Next(def)
					}

					sink(t)
				}
			})
		}
	}
}
