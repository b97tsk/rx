package rx

import (
	"context"
)

// Materialize creates an Observable that represents all of the notifications
// from the source Observable as NEXT emissions, and then completes.
func (Operators) Materialize() Operator {
	return func(source Observable) Observable {
		return func(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
			return source.Subscribe(ctx, func(t Notification) {
				sink.Next(t)
				if !t.HasValue {
					sink.Complete()
				}
			})
		}
	}
}
