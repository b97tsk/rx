package rx

import (
	"context"
)

// Throw creates an Observable that emits no items to the Observer and
// immediately emits an ERROR notification.
func Throw(err error) Observable {
	return func(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
		sink.Error(err)
		return Done(ctx)
	}
}
