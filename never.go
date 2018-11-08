package rx

import (
	"context"
)

var never = Observable{}.Lift(
	func(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
		return context.WithCancel(ctx)
	},
)

// Never creates an Observable that never emits anything.
func Never() Observable {
	return never
}
