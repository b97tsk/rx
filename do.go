package rx

import (
	"context"
)

// Do creates an Observable that mirrors the source Observable, but perform
// a side effect before each emission.
func (Operators) Do(sink Observer) OperatorFunc {
	return func(source Observable) Observable {
		return source.Lift(
			func(ctx context.Context, notify Observer, source Observable) (context.Context, context.CancelFunc) {
				return source.Subscribe(ctx, func(t Notification) {
					sink(t)
					notify(t)
				})
			},
		)
	}
}
