package rx

import (
	"context"
)

// An Observable is a collection of future values. When an Observable is
// subscribed, its values, when available, are emitted to the specified
// Observer.
type Observable func(context.Context, Observer) (context.Context, context.CancelFunc)

// Pipe stitches operators together into a chain, returns the Observable result
// of all of the operators having been called in the order they were passed in.
func (obs Observable) Pipe(operators ...Operator) Observable {
	for _, op := range operators {
		obs = op(obs)
	}
	return obs
}

// Subscribe invokes an execution of an Observable.
func (obs Observable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	if obs == nil {
		obs = Empty()
	}
	return obs(ctx, sink)
}
