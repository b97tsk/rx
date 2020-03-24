package rx

import (
	"context"
)

// Done returns a canceled context and a noop function.
func Done(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	cancel()
	return ctx, cancel
}

// ProjectToObservable type-asserts each value to be an Observable and returns
// it. If type assertion fails, it returns Throw(ErrNotObservable).
func ProjectToObservable(val interface{}, idx int) Observable {
	if obs, ok := val.(Observable); ok {
		return obs
	}
	return Throw(ErrNotObservable)
}
