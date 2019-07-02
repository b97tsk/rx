package rx

import (
	"context"
)

var (
	canceledCtx context.Context
	nothingToDo = func() {}
)

func init() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	canceledCtx = ctx
}

func defaultCompare(v1, v2 interface{}) bool {
	return v1 == v2
}

func defaultKeySelector(val interface{}) interface{} {
	return val
}

// Done returns a canceled context and a function does nothing.
func Done() (context.Context, context.CancelFunc) {
	return canceledCtx, nothingToDo
}

// ProjectToObservable type-casts each value to an Observable and returns it,
// if failed, returns Throw(ErrNotObservable).
func ProjectToObservable(val interface{}, idx int) Observable {
	if obs, ok := val.(Observable); ok {
		return obs
	}
	return Throw(ErrNotObservable)
}
