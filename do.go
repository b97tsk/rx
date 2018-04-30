package rx

import (
	"context"
)

type doOperator struct {
	source Operator
	target Observer
}

func (op doOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	return op.source.Call(ctx, ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			op.target.Next(t.Value)
			ob.Next(t.Value)
		case t.HasError:
			err := t.Value.(error)
			op.target.Error(err)
			ob.Error(err)
		default:
			op.target.Complete()
			ob.Complete()
		}
	}))
}

// Do creates an Observable that mirrors the source Observable, but perform
// a side effect (re-emits each emission to the target Observer) before each
// emission.
func (o Observable) Do(target Observer) Observable {
	op := doOperator{
		source: o.Op,
		target: target,
	}
	return Observable{op}
}
