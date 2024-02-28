package rx

import (
	"context"
	"sync"
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

// ShareOperator is an [Operator] type for [Share].
type ShareOperator[T any] struct {
	opts shareConfig[T]
}

// WithConnector sets Connector option to a given value.
func (op ShareOperator[T]) WithConnector(connector func() Subject[T]) ShareOperator[T] {
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

type shareObservable[T any] struct {
	mu         sync.Mutex
	source     Observable[T]
	connector  func() Subject[T]
	subject    Subject[T]
	connection context.Context
	disconnect CancelFunc
	shareCount int
}

func (obs *shareObservable[T]) Subscribe(c Context, sink Observer[T]) {
	obs.mu.Lock()

	var unlocked bool

	defer func() {
		if !unlocked {
			obs.mu.Unlock()
		}
	}()

	if obs.subject.Observable == nil {
		obs.subject = Try01(obs.connector, func() { sink.Error(ErrOops) })
	}

	c, cancel := c.WithCancel()
	sink = sink.OnLastNotification(cancel)

	obs.subject.Subscribe(c, sink)

	select {
	default:
	case <-c.Done():
		return
	}

	connection := obs.connection

	obs.shareCount++

	if connection == nil {
		w, cancelw := NewBackgroundContext().WithCancel()

		connection = w.Context
		obs.connection = w.Context
		obs.disconnect = cancelw

		sink := obs.subject.Observer

		obs.mu.Unlock()
		unlocked = true

		obs.source.Subscribe(w, func(n Notification[T]) {
			switch n.Kind {
			case KindNext:
				sink(n)
			case KindError, KindComplete:
				cancelw()

				obs.mu.Lock()

				if connection == obs.connection {
					obs.subject = Subject[T]{}
					obs.connection = nil
					obs.disconnect = nil
					obs.shareCount = 0
				}

				obs.mu.Unlock()

				sink(n)
			}
		})
	}

	c.AfterFunc(func() {
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
