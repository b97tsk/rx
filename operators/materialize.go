package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// Materialize represents all of the Notifications from the source as values,
// and then completes.
func Materialize() rx.Operator {
	return materialize
}

func materialize(source rx.Observable) rx.Observable {
	return func(ctx context.Context, sink rx.Observer) {
		source.Subscribe(ctx, func(t rx.Notification) {
			sink.Next(t)

			if !t.HasValue {
				sink.Complete()
			}
		})
	}
}
