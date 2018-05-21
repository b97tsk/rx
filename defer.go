package rx

import (
	"context"
)

type deferOperator struct {
	Func func() Observable
}

func (op deferOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	return op.Func().Subscribe(ctx, sink)
}

// Defer creates an Observable that, on subscribe, calls an Observable
// factory to make an Observable for each new Observer.
//
// Defer creates the Observable lazily, that is, only when it is subscribed.
func Defer(f func() Observable) Observable {
	op := deferOperator{f}
	return Observable{}.Lift(op.Call)
}
