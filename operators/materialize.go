package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

func materialize(source rx.Observable) rx.Observable {
	return func(ctx context.Context, sink rx.Observer) (context.Context, context.CancelFunc) {
		return source.Subscribe(ctx, func(t rx.Notification) {
			sink.Next(t)
			if !t.HasValue {
				sink.Complete()
			}
		})
	}
}

// Materialize creates an Observable that represents all of the notifications
// from the source Observable as NEXT emissions, and then completes.
func Materialize() rx.Operator {
	return materialize
}
