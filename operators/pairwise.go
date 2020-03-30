package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// Pairwise creates an Observable that groups pairs of consecutive emissions
// together and emits them as a slice of two values.
func Pairwise() rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return func(ctx context.Context, sink rx.Observer) (context.Context, context.CancelFunc) {
			var (
				prev    interface{}
				hasPrev bool
			)
			return source.Subscribe(ctx, func(t rx.Notification) {
				switch {
				case t.HasValue:
					if hasPrev {
						sink.Next([]interface{}{prev, t.Value})
						prev = t.Value
					} else {
						prev = t.Value
						hasPrev = true
					}
				default:
					sink(t)
				}
			})
		}
	}
}
