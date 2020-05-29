package operators

import (
	"context"
	"time"

	"github.com/b97tsk/rx"
)

type sampleObservable struct {
	Source   rx.Observable
	Notifier rx.Observable
}

func (obs sampleObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	type X struct {
		LatestValue    interface{}
		HasLatestValue bool
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	obs.Notifier.Subscribe(ctx, func(t rx.Notification) {
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
	})

	if ctx.Err() != nil {
		return
	}

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

// Sample creates an Observable that emits the most recently emitted value from
// the source Observable whenever another Observable, the notifier, emits.
//
// It's like SampleTime, but samples whenever the notifier Observable emits
// something.
func Sample(notifier rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := sampleObservable{source, notifier}
		return rx.Create(obs.Subscribe)
	}
}

// SampleTime creates an Observable that emits the most recently emitted value
// from the source Observable within periodic time intervals.
func SampleTime(period time.Duration) rx.Operator {
	return Sample(rx.Ticker(period))
}
