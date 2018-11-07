package rx

import (
	"context"
)

type goOperator struct{}

func (op goOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	go source.Subscribe(ctx, Finally(sink, cancel))
	return ctx, cancel
}

// Go creates an Observable that asynchronously subscribes the source
// Observable in a goroutine.
func (Operators) Go() OperatorFunc {
	return func(source Observable) Observable {
		op := goOperator{}
		return source.Lift(op.Call)
	}
}
