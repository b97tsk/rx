package rx

import (
	"context"
)

// Empty creates an Observable that emits no items to the Observer and
// immediately emits a COMPLETE notification.
func Empty() Observable {
	return func(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
		sink.Complete()
		return Done(ctx)
	}
}
