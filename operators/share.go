package operators

import (
	"context"
	"sync"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/misc"
)

type shareObservable struct {
	mux            sync.Mutex
	cws            misc.ContextWaitService
	source         rx.Observable
	subjectFactory func() rx.Subject
	subject        rx.Subject
	connection     context.Context
	disconnect     context.CancelFunc
	shareCount     int
}

func (obs *shareObservable) Subscribe(ctx context.Context, sink rx.Observer) (context.Context, context.CancelFunc) {
	obs.mux.Lock()
	defer obs.mux.Unlock()

	if obs.subject.Observer == nil {
		obs.subject = obs.subjectFactory()
	}

	ctx, cancel := obs.subject.Subscribe(ctx, sink)
	if ctx.Err() != nil {
		return ctx, cancel
	}

	connection := obs.connection

	if connection == nil {
		ctx, cancel := context.WithCancel(context.Background())

		connection = ctx
		obs.connection = ctx
		obs.disconnect = cancel

		sink := obs.subject.Observer

		go obs.source.Subscribe(ctx, func(t rx.Notification) {
			if t.HasValue {
				sink(t)
				return
			}

			obs.mux.Lock()
			if connection == obs.connection {
				obs.disconnect()
				obs.subject = rx.Subject{}
				obs.connection = nil
				obs.disconnect = nil
				obs.shareCount = 0
			}
			obs.mux.Unlock()

			sink(t)
		})
	}

	obs.shareCount++

	finalize := func() {
		obs.mux.Lock()
		if connection == obs.connection {
			obs.shareCount--
			if obs.shareCount == 0 {
				obs.disconnect()
				obs.subject = rx.Subject{}
				obs.connection = nil
				obs.disconnect = nil
			}
		}
		obs.mux.Unlock()
	}

	for obs.cws == nil || !obs.cws.Submit(ctx, finalize) {
		obs.cws = misc.NewContextWaitService()
	}

	return ctx, cancel
}

// Share returns a new Observable that multicasts (shares) the original
// Observable. When subscribed multiple times, it guarantees that only one
// subscription is made to the source Observable at the same time. When all
// subscribers have unsubscribed it will unsubscribe from the source Observable.
func Share(subjectFactory func() rx.Subject) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := shareObservable{
			source:         source,
			subjectFactory: subjectFactory,
		}
		return obs.Subscribe
	}
}
