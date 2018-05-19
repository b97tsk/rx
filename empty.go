package rx

import (
	"context"
)

type emptyOperator struct{}

func (op emptyOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	sink.Complete()
	return canceledCtx, doNothing
}

// Empty creates an Observable that emits no items to the Observer and
// immediately emits a Complete notification.
func Empty() Observable {
	op := emptyOperator{}
	return Observable{}.Lift(op.Call)
}
