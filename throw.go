package rx

import (
	"context"
)

// Throw creates an Observable that emits no items to the Observer and
// immediately emits a specified error.
func Throw(e error) Observable {
	return func(ctx context.Context, sink Observer) {
		sink.Error(e)
	}
}
