package rx

import (
	"context"
)

// Throw creates an Observable that emits no items to the Observer and
// immediately emits an ERROR notification.
func Throw(err error) Observable {
	return func(parent context.Context, sink Observer) (context.Context, context.CancelFunc) {
		ctx := NewContext(parent)
		sink.Error(err)
		ctx.Unsubscribe(err)
		return ctx, ctx.Cancel
	}
}
