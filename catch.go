package rx

import (
	"context"
)

type catchOperator struct {
	Selector func(error) Observable
}

func (op catchOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			sink(t)
		case t.HasError:
			obsv := op.Selector(t.Value.(error))
			obsv.Subscribe(ctx, Finally(sink, cancel))
		default:
			sink(t)
			cancel()
		}
	})

	return ctx, cancel
}

// Catch catches errors on the Observable to be handled by returning a new
// Observable.
func (Operators) Catch(selector func(error) Observable) OperatorFunc {
	return func(source Observable) Observable {
		op := catchOperator{selector}
		return source.Lift(op.Call)
	}
}
