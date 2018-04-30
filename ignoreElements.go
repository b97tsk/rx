package rx

import (
	"context"
)

type ignoreElementsOperator struct {
	source Operator
}

func (op ignoreElementsOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	return op.source.Call(ctx, ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
		case t.HasError:
			ob.Error(t.Value.(error))
		default:
			ob.Complete()
		}
	}))
}

// IgnoreElements creates an Observable that ignores all items emitted by the
// source Observable and only passes calls of Complete or Error.
func (o Observable) IgnoreElements() Observable {
	op := ignoreElementsOperator{o.Op}
	return Observable{op}
}
