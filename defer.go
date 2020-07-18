package rx

import (
	"context"
)

// Defer creates an Observable that, on subscribe, calls an Observable
// factory to make an Observable for each new Observer.
//
// Deprecated: Insignificant.
func Defer(create func() Observable) Observable {
	return func(ctx context.Context, sink Observer) {
		create().Subscribe(ctx, sink)
	}
}
