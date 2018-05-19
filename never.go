package rx

import (
	"context"
)

type neverOperator struct{}

func (op neverOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	return context.WithCancel(ctx)
}

// Never creates an Observable that never emits anything.
func Never() Observable {
	op := neverOperator{}
	return Observable{}.Lift(op.Call)
}
