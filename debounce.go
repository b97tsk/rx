package rx

import (
	"context"
)

type debounceOperator struct {
	DurationSelector func(interface{}) Observable
}

func (op debounceOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	scheduleCtx, scheduleCancel := Done()

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

				scheduleCtx, scheduleCancel = context.WithCancel(ctx)

				var observer Observer
				observer = func(t Notification) {
					observer = NopObserver
					scheduleCancel()
					if x, ok := <-cx; ok {
						if t.HasError {
							close(cx)
							sink(t)
							return
						}
						if x.HasLatestValue {
							sink.Next(x.LatestValue)
							x.HasLatestValue = false
						}
						cx <- x
					}
				}

				obs := op.DurationSelector(t.Value)
				obs.Subscribe(scheduleCtx, observer.Notify)

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

// Debounce creates an Observable that emits a value from the source Observable
// only after a particular time span, determined by another Observable, has
// passed without another source emission.
//
// It's like DebounceTime, but the time span of emission silence is determined
// by a second Observable.
func (Operators) Debounce(durationSelector func(interface{}) Observable) OperatorFunc {
	return func(source Observable) Observable {
		op := debounceOperator{durationSelector}
		return source.Lift(op.Call)
	}
}
