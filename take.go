package rx

import (
	"context"
)

type takeOperator struct {
	count int
}

func (op takeOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	var (
		count    = op.count
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
					cancel()
				}
			} else {
				observer = NopObserver
				sink.Complete()
				cancel()
			}
		default:
			sink(t)
			cancel()
		}
	}

	source.Subscribe(ctx, observer.Notify)

	return ctx, cancel
}

// Take creates an Observable that emits only the first count values emitted
// by the source Observable.
//
// Take takes the first count values from the source, then completes.
func (o Observable) Take(count int) Observable {
	op := takeOperator{count}
	return o.Lift(op.Call)
}
