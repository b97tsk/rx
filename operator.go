package rx

import (
	"context"
)

// A Operator is a method on the Observable type. When called, they do not
// change the existing Observable instance. Instead, they return a new
// Observable, whose subscription logic is based on the first Observable.
type Operator interface {
	Call(context.Context, Observer) (context.Context, context.CancelFunc)
}

// OperatorFunc is a helper type that lets you easily create an Operator from
// a function.
type OperatorFunc func(context.Context, Observer) (context.Context, context.CancelFunc)

// Call just calls the underlying function.
func (f OperatorFunc) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	return f(ctx, ob)
}
