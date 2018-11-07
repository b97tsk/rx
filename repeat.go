package rx

import (
	"context"
)

type repeatOperator struct {
	Count int
}

func (op repeatOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	var (
		count          = op.Count
		observer       Observer
		avoidRecursive avoidRecursiveCalls
	)

	subscribe := func() {
		source.Subscribe(ctx, observer)
	}

	observer = func(t Notification) {
		switch {
		case t.HasValue:
			sink(t)
		case t.HasError:
			sink(t)
		default:
			if count == 0 {
				sink(t)
			} else {
				if count > 0 {
					count--
				}
				avoidRecursive.Do(subscribe)
			}
		}
	}

	avoidRecursive.Do(subscribe)

	return ctx, cancel
}

// Repeat creates an Observable that repeats the stream of items emitted by the
// source Observable at most count times.
func (Operators) Repeat(count int) OperatorFunc {
	return func(source Observable) Observable {
		if count == 0 {
			return Empty()
		}
		if count > 0 {
			count--
		}
		op := repeatOperator{count}
		return source.Lift(op.Call)
	}
}
