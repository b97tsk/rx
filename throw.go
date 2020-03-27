package rx

import (
	"context"
)

// Throw creates an Observable that emits no items to the Observer and
// immediately emits an ERROR notification.
func Throw(err error) Observable {
	return Create(func(ctx context.Context, sink Observer) { sink.Error(err) })
}
