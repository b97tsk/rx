package rx

import (
	"context"
)

var canceledCtx context.Context

func init() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	canceledCtx = ctx
}

// Done returns a canceled context and a function does nothing.
func Done() (context.Context, context.CancelFunc) {
	return canceledCtx, func() {}
}

// ProjectToObservable type-asserts each value to be an Observable and returns
// it. If type assertion fails, it returns Throw(ErrNotObservable).
func ProjectToObservable(val interface{}, idx int) Observable {
	if obs, ok := val.(Observable); ok {
		return obs
	}
	return Throw(ErrNotObservable)
}
