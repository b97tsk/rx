package rx

import (
	"context"
)

type takeOperator struct {
	Count int
}

func (op takeOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	var (
		count    = op.Count
		observer Observer
	)

	observer = func(t Notification) {
		switch {
		case t.HasValue:
			if count > 0 {
				count--
				if count > 0 {
					sink(t)
				} else {
					observer = NopObserver
					sink(t)
					sink.Complete()
				}
			} else {
				observer = NopObserver
				sink.Complete()
			}
		default:
			sink(t)
		}
	}

	source.Subscribe(ctx, observer.Notify)

	return ctx, cancel
}

// Take creates an Observable that emits only the first count values emitted
// by the source Observable.
//
// Take takes the first count values from the source, then completes.
func (Operators) Take(count int) OperatorFunc {
	return func(source Observable) Observable {
		op := takeOperator{count}
		return source.Lift(op.Call)
	}
}
