package rx

import (
	"context"
)

type throwOperator struct {
	err error
}

func (op throwOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	sink.Error(op.err)
	return canceledCtx, doNothing
}

// Throw creates an Observable that emits no items to the Observer and
// immediately emits an Error notification.
func Throw(err error) Observable {
	op := throwOperator{err}
	return Observable{}.Lift(op.Call)
}
