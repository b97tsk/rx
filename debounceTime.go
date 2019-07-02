package rx

import (
	"context"
	"time"
)

type debounceTimeOperator struct {
	Duration time.Duration
}

func (op debounceTimeOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	_, scheduleCancel := Done()

	sink = Finally(sink, cancel)

	type X struct {
		LatestValue    interface{}
		HasLatestValue bool
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	source.Subscribe(ctx, func(t Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				x.LatestValue = t.Value
				x.HasLatestValue = true

				cx <- x

				scheduleCancel()

				_, scheduleCancel = scheduleOnce(ctx, op.Duration, func() {
					if x, ok := <-cx; ok {
						if x.HasLatestValue {
							sink.Next(x.LatestValue)
							x.HasLatestValue = false
						}
						cx <- x
					}
				})

			case t.HasError:
				close(cx)
				sink(t)

			default:
				close(cx)
				if x.HasLatestValue {
					sink.Next(x.LatestValue)
				}
				sink(t)
			}
		}
	})

	return ctx, cancel
}

// DebounceTime creates an Observable that emits a value from the source
// Observable only after a particular time span has passed without another
// source emission.
//
// It's like Delay, but passes only the most recent value from each burst of
// emissions.
func (Operators) DebounceTime(duration time.Duration) OperatorFunc {
	return func(source Observable) Observable {
		op := debounceTimeOperator{duration}
		return source.Lift(op.Call)
	}
}
