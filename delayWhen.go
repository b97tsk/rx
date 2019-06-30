package rx

import (
	"context"
)

type delayWhenOperator struct {
	DurationSelector func(interface{}, int) Observable
}

func (op delayWhenOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	type X struct {
		Index       int
		ActiveCount int
	}
	cx := make(chan *X, 1)
	cx <- &X{ActiveCount: 1}

	source.Subscribe(ctx, func(t Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				outerIndex := x.Index
				outerValue := t.Value
				x.Index++
				x.ActiveCount++

				cx <- x

				scheduleCtx, scheduleCancel := context.WithCancel(ctx)

				var observer Observer
				observer = func(t Notification) {
					observer = NopObserver
					scheduleCancel()
					if x, ok := <-cx; ok {
						x.ActiveCount--
						switch {
						case t.HasValue:
							sink.Next(outerValue)
							if x.ActiveCount == 0 {
								close(cx)
								sink.Complete()
								return
							}
							cx <- x
						case t.HasError:
							close(cx)
							sink(t)
						default:
							if x.ActiveCount == 0 {
								close(cx)
								sink(t)
								return
							}
							cx <- x
						}
					}
				}

				obs := op.DurationSelector(outerValue, outerIndex)
				obs.Subscribe(scheduleCtx, observer.Notify)

			case t.HasError:
				close(cx)
				sink(t)

			default:
				x.ActiveCount--
				if x.ActiveCount == 0 {
					close(cx)
					sink(t)
					return
				}
				cx <- x
			}
		}
	})

	return ctx, cancel
}

// DelayWhen creates an Observable that delays the emission of items from
// the source Observable by a given time span determined by the emissions of
// another Observable.
//
// It's like Delay, but the time span of the delay duration is determined by
// a second Observable.
func (Operators) DelayWhen(durationSelector func(interface{}, int) Observable) OperatorFunc {
	return func(source Observable) Observable {
		op := delayWhenOperator{durationSelector}
		return source.Lift(op.Call)
	}
}
