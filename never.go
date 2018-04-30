package rx

import (
	"context"
)

type neverOperator struct{}

func (op neverOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	return ctx, noopFunc
}

// Never creates an Observable that never emits anything.
func Never() Observable {
	op := neverOperator{}
	return Observable{op}
}
