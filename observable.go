package rx

import (
	"context"
)

// An Observable is a collection of future values. When an Observable is
// subscribed, its values, when available, are emitted to the specified
// Observer.
//
// It's easy to create your own Observable, but you're responsible to follow
// the Observable Contract that no more emissions pass to the specified
// Observer after an error or a completion has passed to it. Violations of
// this contract result in undefined behaviors.
type Observable func(context.Context, Observer)

// Pipe stitches operators together into a chain, returns the Observable result
// of all of the operators having been called in the order they were passed in.
func (obs Observable) Pipe(operators ...Operator) Observable {
	for _, op := range operators {
		obs = op(obs)
	}
	return obs
}

// Subscribe invokes an execution of an Observable.
func (obs Observable) Subscribe(ctx context.Context, sink Observer) {
	obs(ctx, sink)
}
