package rx

import (
	"context"
)

// An Observable is a collection of future values. When an Observable is
// subscribed, its values, when available, are emitted to the specified
// Observer.
type Observable struct {
	*observable
}

type observable struct {
	source Observable
	op     Operator
}

// Lift creates a new Observable, with this Observable as the source, and
// the passed Operator defined as the new Observable's Operator.
func (o Observable) Lift(op Operator) Observable {
	return Observable{&observable{o, op}}
}

// Pipe stitches Operators together into a chain, returns the Observable result
// of all of the Operators having been called in the order they were passed in.
func (o Observable) Pipe(operations ...Operator) Observable {
	for _, op := range operations {
		o = o.Lift(op)
	}
	return o
}

// Subscribe invokes an execution of an Observable.
func (o Observable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	return o.op(ctx, sink, o.source)
}
