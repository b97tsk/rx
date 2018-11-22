package rx

import (
	"context"
)

// Defer creates an Observable that, on subscribe, calls an Observable
// factory to make an Observable for each new Observer.
//
// Defer creates the Observable lazily, that is, only when it is subscribed.
func Defer(create func() Observable) Observable {
	return Observable{}.Lift(
		func(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
			return create().Subscribe(ctx, sink)
		},
	)
}
