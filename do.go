package rx

import (
	"context"
)

type doOperator struct {
	target Observer
}

func (op doOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	return source.Subscribe(ctx, func(t Notification) {
		t.Observe(op.target)
		t.Observe(ob)
	})
}

// Do creates an Observable that mirrors the source Observable, but perform
// a side effect before each emission.
func (o Observable) Do(target Observer) Observable {
	op := doOperator{target}
	return o.Lift(op.Call)
}
