package rx

import (
	"context"
)

type ignoreElementsOperator struct{}

func (op ignoreElementsOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	return source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
		default:
			t.Observe(ob)
		}
	})
}

// IgnoreElements creates an Observable that ignores all items emitted by the
// source Observable and only passes calls of Complete or Error.
func (o Observable) IgnoreElements() Observable {
	op := ignoreElementsOperator{}
	return o.Lift(op.Call)
}
