package rx

import (
	"context"
	"time"
)

type subscribeOnObservable struct {
	Source   Observable
	Duration time.Duration
}

func (obs subscribeOnObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	scheduleOnce(ctx, obs.Duration, func() {
		obs.Source.Subscribe(ctx, sink)
	})

	return ctx, cancel
}

// SubscribeOn creates an Observable that asynchronously subscribes to the
// source Observable after waits for the duration to elapse.
func (Operators) SubscribeOn(d time.Duration) Operator {
	return func(source Observable) Observable {
		return subscribeOnObservable{source, d}.Subscribe
	}
}
