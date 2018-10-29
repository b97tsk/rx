package rx

import (
	"context"
)

type firstOperator struct{}

func (op firstOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	var observer Observer
	observer = func(t Notification) {
		switch {
		case t.HasValue:
			observer = NopObserver
			sink(t)
			sink.Complete()
			cancel()
		case t.HasError:
			sink(t)
			cancel()
		default:
			sink.Error(ErrEmpty)
			cancel()
		}
	}
	source.Subscribe(ctx, observer.Notify)

	return ctx, cancel
}

// First creates an Observable that emits only the first value (or the first
// value that meets some condition) emitted by the source Observable.
func (Operators) First() OperatorFunc {
	return func(source Observable) Observable {
		op := firstOperator{}
		return source.Lift(op.Call)
	}
}
