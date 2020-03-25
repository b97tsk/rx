package rx

import (
	"context"
)

// Empty creates an Observable that emits no items to the Observer and
// immediately emits a COMPLETE notification.
func Empty() Observable {
	return func(parent context.Context, sink Observer) (context.Context, context.CancelFunc) {
		ctx := NewContext(parent)
		sink.Complete()
		ctx.Unsubscribe(Complete)
		return ctx, ctx.Cancel
	}
}
