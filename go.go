package rx

import (
	"context"
)

// Go creates an Observable that mirrors the source Observable in a goroutine.
func (Operators) Go() Operator {
	return func(source Observable) Observable {
		return func(parent context.Context, sink Observer) (context.Context, context.CancelFunc) {
			ctx := NewContext(parent)
			go source.Subscribe(ctx, DoAtLast(sink, ctx.AtLast))
			return ctx, ctx.Cancel
		}
	}
}
