package rx

import (
	"context"
)

type sampleObservable struct {
	Source   Observable
	Notifier Observable
}

func (obs sampleObservable) Subscribe(ctx context.Context, sink Observer) {
	type X struct {
		LatestValue    interface{}
		HasLatestValue bool
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	obs.Notifier.Subscribe(ctx, func(t Notification) {
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
}

// Sample creates an Observable that emits the most recently emitted value from
// the source Observable whenever another Observable, the notifier, emits.
//
// It's like SampleTime, but samples whenever the notifier Observable emits
// something.
func (Operators) Sample(notifier Observable) Operator {
	return func(source Observable) Observable {
		obs := sampleObservable{source, notifier}
		return Create(obs.Subscribe)
	}
}
