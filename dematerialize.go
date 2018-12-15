package rx

import (
	"context"
)

type dematerializeOperator struct{}

func (op dematerializeOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	var observer Observer
	observer = func(t Notification) {
		switch {
		case t.HasValue:
			if t, ok := t.Value.(Notification); ok {
				switch {
				case t.HasValue:
					sink(t)
				default:
					observer = NopObserver
					sink(t)
				}
			} else {
				observer = NopObserver
				sink.Error(ErrNotNotification)
			}
		default:
			sink(t)
		}
	}
	source.Subscribe(ctx, observer.Notify)

	return ctx, cancel
}

// Dematerialize converts an Observable of Notification objects into the
// emissions that they represent. It's the opposite of Materialize.
func (Operators) Dematerialize() OperatorFunc {
	return func(source Observable) Observable {
		op := dematerializeOperator{}
		return source.Lift(op.Call)
	}
}
