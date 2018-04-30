package rx

import (
	"context"
)

// An Observable is a collection of future values. When an Observable is
// subscribed, its values, when available, are emitted to the specified
// Observer.
type Observable struct {
	Op Operator
}

// Subscribe invokes an execution of an Observable.
func (o Observable) Subscribe(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	return o.Op.Call(ctx, ob)
}

type operatorWithOptions interface {
	ApplyOptions([]Option) Operator
}

// With applies Options to the source Observable, and returns a new one.
func (o Observable) With(options []Option) Observable {
	if len(options) > 0 {
		if op, ok := o.Op.(operatorWithOptions); ok {
			o.Op = op.ApplyOptions(options)
		} else {
			panic(ErrNoOptions)
		}
	}
	return o
}
