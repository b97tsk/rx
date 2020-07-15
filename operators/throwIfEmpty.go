package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// ThrowIfEmpty creates an Observable that mirrors the source Observable, if
// the source turns out to be empty, throws a specified error.
func ThrowIfEmpty(e error) rx.Operator {
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
					if hasValue {
						sink(t)
						return
					}
					sink.Error(e)
				}
			})
		}
	}
}
