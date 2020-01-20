package rx

import (
	"context"
)

type delayWhenObservable struct {
	Source           Observable
	DurationSelector func(interface{}, int) Observable
}

func (obs delayWhenObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	type X struct {
		Index       int
		ActiveCount int
	}
	cx := make(chan *X, 1)
	cx <- &X{ActiveCount: 1}

	obs.Source.Subscribe(ctx, func(t Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				sourceIndex := x.Index
				sourceValue := t.Value
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
							sink.Next(sourceValue)
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

				obs := obs.DurationSelector(sourceValue, sourceIndex)
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
func (Operators) DelayWhen(durationSelector func(interface{}, int) Observable) Operator {
	return func(source Observable) Observable {
		return delayWhenObservable{source, durationSelector}.Subscribe
	}
}
