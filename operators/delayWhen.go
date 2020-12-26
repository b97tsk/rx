package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/atomic"
)

// DelayWhen creates an Observable that delays the emission of items from
// the source Observable by a given time span determined by the emissions of
// another Observable.
//
// It's like Delay, but the time span of the delay duration is determined by
// a second Observable.
func DelayWhen(durationSelector func(interface{}, int) rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return delayWhenObservable{source, durationSelector}.Subscribe
	}
}

type delayWhenObservable struct {
	Source           rx.Observable
	DurationSelector func(interface{}, int) rx.Observable
}

func (obs delayWhenObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel).Mutex()

	sourceIndex := -1
	workers := atomic.Uint32(1)

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		switch {
		case t.HasValue:
			sourceIndex++

			sourceValue := t.Value

			workers.Add(1)

			scheduleCtx, scheduleCancel := context.WithCancel(ctx)

			var observer rx.Observer

			observer = func(t rx.Notification) {
				observer = rx.Noop

				scheduleCancel()

				switch {
				case t.HasValue:
					sink.Next(sourceValue)

					if workers.Sub(1) == 0 {
						sink.Complete()
					}

				case t.HasError:
					sink(t)

				default:
					if workers.Sub(1) == 0 {
						sink(t)
					}
				}
			}

			obs1 := obs.DurationSelector(sourceValue, sourceIndex)

			obs1.Subscribe(scheduleCtx, observer.Sink)

		case t.HasError:
			sink(t)

		default:
			if workers.Sub(1) == 0 {
				sink(t)
			}
		}
	})
}
