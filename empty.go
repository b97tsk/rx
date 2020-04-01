package rx

import (
	"context"
)

var empty = Create(func(ctx context.Context, sink Observer) { sink.Complete() })

// Empty returns an Observable that emits no items to the Observer and
// immediately emits a COMPLETE notification.
func Empty() Observable {
	return empty
}
