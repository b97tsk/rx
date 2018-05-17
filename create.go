package rx

import (
	"context"
)

type createOperator struct {
	f func(context.Context, Observer) (context.Context, context.CancelFunc)
}

func (op createOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	return op.f(ctx, ob)
}

// Create creates a new Observable, that will execute the specified function
// when an Observer subscribes to it.
//
// Create custom Observable, that does whatever you like.
func Create(f func(context.Context, Observer) (context.Context, context.CancelFunc)) Observable {
	op := createOperator{f}
	return Observable{}.Lift(op.Call)
}
