package rx

import (
	"context"
)

type retryOperator struct {
	Count int
}

func (op retryOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	var (
		count    = op.Count
		observer Observer
	)

	observer = func(t Notification) {
		switch {
		case t.HasValue:
			sink(t)
		case t.HasError:
			if count == 0 {
				sink(t)
				cancel()
			} else {
				if count > 0 {
					count--
				}
				source.Subscribe(ctx, observer)
			}
		default:
			sink(t)
			cancel()
		}
	}

	source.Subscribe(ctx, observer)

	return ctx, cancel
}

// Retry creates an Observable that mirrors the source Observable with the
// exception of an Error. If the source Observable calls Error, this method
// will resubscribe to the source Observable for a maximum of count
// resubscriptions rather than propagating the Error call.
func (o Observable) Retry(count int) Observable {
	if count == 0 {
		return o
	}
	op := retryOperator{count}
	return o.Lift(op.Call)
}
