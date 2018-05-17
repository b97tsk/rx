package rx

import (
	"context"
)

type catchOperator struct {
	selector func(error) Observable
}

func (op catchOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			t.Observe(ob)
		case t.HasError:
			obsv := op.selector(t.Value.(error))
			obsv.Subscribe(ctx, withFinalizer(ob, cancel))
		default:
			t.Observe(ob)
			cancel()
		}
	})

	return ctx, cancel
}

// Catch catches errors on the Observable to be handled by returning a new
// Observable.
func (o Observable) Catch(selector func(error) Observable) Observable {
	op := catchOperator{selector}
	return o.Lift(op.Call)
}
