package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// DefaultIfEmpty mirrors the source, emits a given value if the source
// completes without emitting any value.
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
