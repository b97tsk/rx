package rx

import (
	"context"
)

type ignoreElementsOperator struct{}

func (op ignoreElementsOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	return source.Subscribe(ctx, func(t Notification) {
		if t.HasValue {
			return
		}
		sink(t)
	})
}

// IgnoreElements creates an Observable that ignores all values emitted by the
// source Observable and only passes Complete or Error emission.
func (Operators) IgnoreElements() OperatorFunc {
	return func(source Observable) Observable {
		op := ignoreElementsOperator{}
		return source.Lift(op.Call)
	}
}
