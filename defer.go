package rx

import (
	"context"
)

type deferOperator struct {
	factory func() Observable
}

func (op deferOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	return op.factory().Subscribe(ctx, ob)
}

// Defer creates an Observable that, on subscribe, calls an Observable
// factory to make an Observable for each new Observer.
//
// Defer creates the Observable lazily, that is, only when it is subscribed.
func Defer(f func() Observable) Observable {
	op := deferOperator{f}
	return Observable{}.Lift(op.Call)
}
