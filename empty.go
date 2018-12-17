package rx

import (
	"context"
)

var empty = Observable{}.Lift(
	func(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
		sink.Complete()
		return Done()
	},
)

// Empty creates an Observable that emits no items to the Observer and
// immediately emits a COMPLETE notification.
func Empty() Observable {
	return empty
}
