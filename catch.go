package rx

import (
	"context"
)

type catchOperator struct {
	source   Operator
	selector func(error) Observable
}

func (op catchOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	op.source.Call(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			ob.Next(t.Value)
		case t.HasError:
			obsv := op.selector(t.Value.(error))
			obsv.Subscribe(ctx, withFinalizer(ob, cancel))
		default:
			ob.Complete()
			cancel()
		}
	})

	return ctx, cancel
}

// Catch catches errors on the Observable to be handled by returning a new
// Observable.
func (o Observable) Catch(selector func(error) Observable) Observable {
	op := catchOperator{
		source:   o.Op,
		selector: selector,
	}
	return Observable{op}
}
