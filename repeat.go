package rx

import (
	"context"
)

type repeatOperator struct {
	source Operator
	count  int
}

func (op repeatOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	count := op.count

	var observer Observer

	observer = ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			ob.Next(t.Value)
		case t.HasError:
			ob.Error(t.Value.(error))
			cancel()
		default:
			if count == 0 {
				ob.Complete()
				cancel()
			} else {
				if count > 0 {
					count--
				}
				op.source.Call(ctx, observer)
			}
		}
	})

	op.source.Call(ctx, observer)

	return ctx, cancel
}

// Repeat creates an Observable that repeats the stream of items emitted by the
// source Observable at most count times.
func (o Observable) Repeat(count int) Observable {
	if count == 0 {
		return Empty()
	}
	if count > 0 {
		count--
	}
	op := repeatOperator{
		source: o.Op,
		count:  count,
	}
	return Observable{op}
}
