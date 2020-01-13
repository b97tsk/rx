package rx

import (
	"context"
)

// Go creates an Observable that mirrors the source Observable in a goroutine.
func (Operators) Go() Operator {
	return func(source Observable) Observable {
		return func(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
			ctx, cancel := context.WithCancel(ctx)
			go source.Subscribe(ctx, Finally(sink, cancel))
			return ctx, cancel
		}
	}
}
