package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

type delayWhenObservable struct {
	Source           rx.Observable
	DurationSelector func(interface{}, int) rx.Observable
}

func (obs delayWhenObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	type X struct {
		Index       int
		ActiveCount int
	}
	cx := make(chan *X, 1)
	cx <- &X{ActiveCount: 1}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				sourceIndex := x.Index
				sourceValue := t.Value
				x.Index++
				x.ActiveCount++

				cx <- x

				scheduleCtx, scheduleCancel := context.WithCancel(ctx)

				var observer rx.Observer
				observer = func(t rx.Notification) {
					observer = rx.NopObserver
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
}

// DelayWhen creates an Observable that delays the emission of items from
// the source Observable by a given time span determined by the emissions of
// another Observable.
//
// It's like Delay, but the time span of the delay duration is determined by
// a second Observable.
func DelayWhen(durationSelector func(interface{}, int) rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := delayWhenObservable{source, durationSelector}
		return rx.Create(obs.Subscribe)
	}
}
