package rx

import (
	"context"
	"sync"

	"github.com/b97tsk/rx/internal/ctxwatch"
	"github.com/b97tsk/rx/internal/waitgroup"
)

// Share returns a new Observable that multicasts (shares) the source
// Observable. When subscribed multiple times, it guarantees that only one
// subscription is made to the source at the same time. When all subscribers
// have unsubscribed it will unsubscribe from the source.
func Share[T any]() ShareOperator[T] {
	return ShareOperator[T]{
		opts: shareConfig[T]{
			Connector: Multicast[T],
		},
	}
}

type shareConfig[T any] struct {
	Connector func() Subject[T]
}

// ShareOperator is an Operator type for Share.
type ShareOperator[T any] struct {
	opts shareConfig[T]
}

// WithConnector sets Connector option to a given value.
func (op ShareOperator[T]) WithConnector(connector func() Subject[T]) ShareOperator[T] {
	if connector == nil {
		panic("connector == nil")
	}

	op.opts.Connector = connector

	return op
}

// Apply implements the Operator interface.
func (op ShareOperator[T]) Apply(source Observable[T]) Observable[T] {
	obs := shareObservable[T]{
		source:    source,
		connector: op.opts.Connector,
	}

	return obs.Subscribe
}

// AsOperator converts op to an Operator.
//
// Once type inference has improved in Go, this method will be removed.
func (op ShareOperator[T]) AsOperator() Operator[T, T] { return op }

type shareObservable[T any] struct {
	mu         sync.Mutex
	source     Observable[T]
	connector  func() Subject[T]
	subject    Subject[T]
	connection context.Context
	disconnect context.CancelFunc
	shareCount int
}

func (obs *shareObservable[T]) Subscribe(ctx context.Context, sink Observer[T]) {
	obs.mu.Lock()

	var unlocked bool

	defer func() {
		if !unlocked {
			obs.mu.Unlock()
		}
	}()

	if obs.subject.Observable == nil {
		obs.subject = obs.connector()
	}

	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel)

	obs.subject.Subscribe(ctx, sink)

	if getErr(ctx) != nil {
		return
	}

	connection := obs.connection

	obs.shareCount++

	if connection == nil {
		ctx, cancel := context.WithCancel(context.Background())

		connection = ctx
		obs.connection = ctx
		obs.disconnect = cancel

		sink := obs.subject.Observer

		obs.mu.Unlock()

		unlocked = true

		obs.source.Subscribe(ctx, func(n Notification[T]) {
			if n.HasValue {
				sink(n)
				return
			}

			cancel()

			obs.mu.Lock()

			if connection == obs.connection {
				obs.subject = Subject[T]{}
				obs.connection = nil
				obs.disconnect = nil
				obs.shareCount = 0
			}

			obs.mu.Unlock()

			sink(n)
		})
	}

	wg := waitgroup.Get(ctx)
	if wg != nil {
		wg.Add(1)
	}

	ctxwatch.Add(ctx, func() {
		if wg != nil {
			defer wg.Done()
		}

		obs.mu.Lock()

		if connection == obs.connection {
			obs.shareCount--

			if obs.shareCount == 0 {
				obs.disconnect()

				obs.subject = Subject[T]{}
				obs.connection = nil
				obs.disconnect = nil
			}
		}

		obs.mu.Unlock()
	})
}
