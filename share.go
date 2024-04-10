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
		ts: shareConfig[T]{
			Connector: Multicast[T],
		},
	}
}

type shareConfig[T any] struct {
	Connector func() Subject[T]
}

// ShareOperator is an [Operator] type for [Share].
type ShareOperator[T any] struct {
	ts shareConfig[T]
}

// WithConnector sets Connector option to a given value.
func (op ShareOperator[T]) WithConnector(connector func() Subject[T]) ShareOperator[T] {
	op.ts.Connector = connector
	return op
}

// Apply implements the Operator interface.
func (op ShareOperator[T]) Apply(source Observable[T]) Observable[T] {
	obs := shareObservable[T]{
		Source:    source,
		Connector: op.ts.Connector,
	}

	return obs.Subscribe
}

type shareObservable[T any] struct {
	Mutex      sync.Mutex
	Source     Observable[T]
	Connector  func() Subject[T]
	Subject    Subject[T]
	Connection context.Context
	Disconnect CancelFunc
	ShareCount int
}

func (obs *shareObservable[T]) Subscribe(c Context, o Observer[T]) {
	obs.Mutex.Lock()

	var unlocked bool

	defer func() {
		if !unlocked {
			obs.Mutex.Unlock()
		}
	}()

	if obs.Subject.Observable == nil {
		obs.Subject = Try01(obs.Connector, func() { o.Error(ErrOops) })
	}

	c, cancel := c.WithCancel()
	o = o.DoOnTermination(cancel)

	obs.Subject.Subscribe(c, o)

	select {
	default:
	case <-c.Done():
		return
	}

	connection := obs.Connection

	obs.ShareCount++

	if connection == nil {
		w, cancelw := NewBackgroundContext().WithCancel()

		connection = w.Context
		obs.Connection = w.Context
		obs.Disconnect = cancelw

		o := obs.Subject.Observer

		obs.Mutex.Unlock()
		unlocked = true

		obs.Source.Subscribe(w, func(n Notification[T]) {
			switch n.Kind {
			case KindNext:
				o.Emit(n)
			case KindError, KindComplete:
				cancelw()

				obs.Mutex.Lock()

				if connection == obs.Connection {
					obs.Subject = Subject[T]{}
					obs.Connection = nil
					obs.Disconnect = nil
					obs.ShareCount = 0
				}

				obs.Mutex.Unlock()

				o.Emit(n)
			}
		})
	}

	c.AfterFunc(func() {
		obs.Mutex.Lock()

		if connection == obs.Connection {
			obs.ShareCount--

			if obs.ShareCount == 0 {
				obs.Disconnect()

				obs.Subject = Subject[T]{}
				obs.Connection = nil
				obs.Disconnect = nil
			}
		}

		obs.Mutex.Unlock()
	})
}
