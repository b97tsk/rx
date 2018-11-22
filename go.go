package rx

import (
	"context"
)

// Go creates an Observable that asynchronously subscribes the source
// Observable in a goroutine.
func (Operators) Go() OperatorFunc {
	return func(source Observable) Observable {
		return source.Lift(
			func(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(ctx)
				go source.Subscribe(ctx, Finally(sink, cancel))
				return ctx, cancel
			},
		)
	}
}
