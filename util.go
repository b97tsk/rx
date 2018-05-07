package rx

import (
	"context"
)

var (
	canceledCtx context.Context
	noopFunc    = func() {}
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

func projectToObservable(val interface{}, index int) Observable {
	if obsv, ok := val.(Observable); ok {
		return obsv
	}
	return Throw(ErrNotObservable)
}
