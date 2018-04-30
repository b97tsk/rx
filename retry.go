package rx

import (
	"context"
)

type retryOperator struct {
	source Operator
	count  int
}

func (op retryOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	count := op.count

	var observer Observer

	observer = ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			ob.Next(t.Value)
		case t.HasError:
			if count == 0 {
				ob.Error(t.Value.(error))
				cancel()
			} else {
				if count > 0 {
					count--
				}
				op.source.Call(ctx, observer)
			}
		default:
			ob.Complete()
			cancel()
		}
	})

	op.source.Call(ctx, observer)

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
	op := retryOperator{
		source: o.Op,
		count:  count,
	}
	return Observable{op}
}
