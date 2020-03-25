package rx

import (
	"context"
	"time"
)

type subscribeOnObservable struct {
	Source   Observable
	Duration time.Duration
}

func (obs subscribeOnObservable) Subscribe(parent context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx := NewContext(parent)

	sink = DoAtLast(sink, ctx.AtLast)

	scheduleOnce(ctx, obs.Duration, func() {
		obs.Source.Subscribe(ctx, sink)
	})

	return ctx, ctx.Cancel
}

// SubscribeOn creates an Observable that asynchronously subscribes to the
// source Observable after waits for the duration to elapse.
func (Operators) SubscribeOn(d time.Duration) Operator {
	return func(source Observable) Observable {
		return subscribeOnObservable{source, d}.Subscribe
	}
}
