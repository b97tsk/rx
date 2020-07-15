package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// Catch catches errors on the Observable to be handled by returning a new
// Observable.
func Catch(selector func(error) rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return func(ctx context.Context, sink rx.Observer) {
			source.Subscribe(ctx, func(t rx.Notification) {
				switch {
				case t.HasValue:
					sink(t)
				case t.HasError:
					obs := selector(t.Error)
					obs.Subscribe(ctx, sink)
				default:
					sink(t)
				}
			})
		}
	}
}
