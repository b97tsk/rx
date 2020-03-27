package rx

import (
	"context"
)

// Empty creates an Observable that emits no items to the Observer and
// immediately emits a COMPLETE notification.
func Empty() Observable {
	return Create(func(ctx context.Context, sink Observer) { sink.Complete() })
}
