package operators

import (
	"context"
	"time"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/x/schedule"
)

type sampleTimeObservable struct {
	Source   rx.Observable
	Duration time.Duration
}

func (obs sampleTimeObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	type X struct {
		LatestValue    interface{}
		HasLatestValue bool
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	schedule.Schedule(ctx, obs.Duration, func() {
		if x, ok := <-cx; ok {
			if x.HasLatestValue {
				sink.Next(x.LatestValue)
				x.HasLatestValue = false
			}
			cx <- x
		}
	})

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
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
}

// SampleTime creates an Observable that emits the most recently emitted value
// from the source Observable within periodic time intervals.
func SampleTime(interval time.Duration) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := sampleTimeObservable{source, interval}
		return rx.Create(obs.Subscribe)
	}
}
