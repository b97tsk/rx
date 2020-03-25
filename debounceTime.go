package rx

import (
	"context"
	"time"
)

type debounceTimeObservable struct {
	Source   Observable
	Duration time.Duration
}

func (obs debounceTimeObservable) Subscribe(parent context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx := NewContext(parent)

	sink = DoAtLast(sink, ctx.AtLast)

	type X struct {
		LatestValue    interface{}
		HasLatestValue bool
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	var scheduleCancel context.CancelFunc

	obs.Source.Subscribe(ctx, func(t Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				x.LatestValue = t.Value
				x.HasLatestValue = true

				cx <- x

				if scheduleCancel != nil {
					scheduleCancel()
				}

				_, scheduleCancel = scheduleOnce(ctx, obs.Duration, func() {
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

	return ctx, ctx.Cancel
}

// DebounceTime creates an Observable that emits a value from the source
// Observable only after a particular time span has passed without another
// source emission.
//
// It's like Delay, but passes only the most recent value from each burst of
// emissions.
func (Operators) DebounceTime(duration time.Duration) Operator {
	return func(source Observable) Observable {
		return debounceTimeObservable{source, duration}.Subscribe
	}
}
