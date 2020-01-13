package rx

import (
	"context"
	"time"
)

type sampleTimeObservable struct {
	Source   Observable
	Duration time.Duration
}

func (obs sampleTimeObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	type X struct {
		LatestValue    interface{}
		HasLatestValue bool
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	schedule(ctx, obs.Duration, func() {
		if x, ok := <-cx; ok {
			if x.HasLatestValue {
				sink.Next(x.LatestValue)
				x.HasLatestValue = false
			}
			cx <- x
		}
	})

	obs.Source.Subscribe(ctx, func(t Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				x.LatestValue = t.Value
				x.HasLatestValue = true
				cx <- x
			default:
				close(cx)
				sink(t)
			}
		}
	})

	return ctx, cancel
}

// SampleTime creates an Observable that emits the most recently emitted value
// from the source Observable within periodic time intervals.
func (Operators) SampleTime(interval time.Duration) Operator {
	return func(source Observable) Observable {
		return sampleTimeObservable{source, interval}.Subscribe
	}
}
