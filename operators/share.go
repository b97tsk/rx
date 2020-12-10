package operators

import (
	"context"
	"sync"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/ctxutil"
)

// Share returns a new Observable that multicasts (shares) the original
// Observable. When subscribed multiple times, it guarantees that only one
// subscription is made to the source Observable at the same time. When all
// subscribers have unsubscribed it will unsubscribe from the source Observable.
func Share(doubleFactory rx.DoubleFactory) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := shareObservable{
			source:        source,
			doubleFactory: doubleFactory,
		}

		return obs.Subscribe
	}
}

type shareObservable struct {
	mu            sync.Mutex
	cws           ctxutil.ContextWaitService
	source        rx.Observable
	doubleFactory rx.DoubleFactory
	double        rx.Double
	connection    context.Context
	disconnect    context.CancelFunc
	shareCount    int
}

func (obs *shareObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	obs.mu.Lock()
	defer obs.mu.Unlock()

	if obs.double.Observable == nil {
		obs.double = obs.doubleFactory()
	}

	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel)

	obs.double.Subscribe(ctx, sink)

	if ctx.Err() != nil {
		return
	}

	connection := obs.connection

	if connection == nil {
		ctx, cancel := context.WithCancel(context.Background())

		connection = ctx
		obs.connection = ctx
		obs.disconnect = cancel

		sink := obs.double.Observer

		go obs.source.Subscribe(ctx, func(t rx.Notification) {
			if t.HasValue {
				sink(t)
				return
			}

			cancel()

			obs.mu.Lock()

			if connection == obs.connection {
				obs.double = rx.Double{}
				obs.connection = nil
				obs.disconnect = nil
				obs.shareCount = 0
			}

			obs.mu.Unlock()

			sink(t)
		})
	}

	obs.shareCount++

	finalize := func() {
		obs.mu.Lock()

		if connection == obs.connection {
			obs.shareCount--

			if obs.shareCount == 0 {
				obs.disconnect()
				obs.double = rx.Double{}
				obs.connection = nil
				obs.disconnect = nil
			}
		}

		obs.mu.Unlock()
	}

	for obs.cws == nil || !obs.cws.Submit(ctx, finalize) {
		obs.cws = ctxutil.NewContextWaitService()
	}
}
