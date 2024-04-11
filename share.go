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
	ob := shareObservable[T]{
		Source:    source,
		Connector: op.ts.Connector,
	}

	return ob.Subscribe
}

type shareObservable[T any] struct {
	Mu         sync.Mutex
	Source     Observable[T]
	Connector  func() Subject[T]
	Subject    Subject[T]
	Connection context.Context
	Disconnect CancelFunc
	ShareCount int
}

func (ob *shareObservable[T]) Subscribe(c Context, o Observer[T]) {
	ob.Mu.Lock()

	var unlocked bool

	defer func() {
		if !unlocked {
			ob.Mu.Unlock()
		}
	}()

	if ob.Subject.Observable == nil {
		ob.Subject = Try01(ob.Connector, func() { o.Error(ErrOops) })
	}

	c, cancel := c.WithCancel()
	o = o.DoOnTermination(cancel)

	ob.Subject.Subscribe(c, o)

	select {
	default:
	case <-c.Done():
		return
	}

	connection := ob.Connection

	ob.ShareCount++

	if connection == nil {
		w, cancelw := NewBackgroundContext().WithCancel()

		connection = w.Context
		ob.Connection = w.Context
		ob.Disconnect = cancelw

		o := ob.Subject.Observer

		ob.Mu.Unlock()
		unlocked = true

		ob.Source.Subscribe(w, func(n Notification[T]) {
			switch n.Kind {
			case KindNext:
				o.Emit(n)
			case KindError, KindComplete:
				cancelw()

				ob.Mu.Lock()

				if connection == ob.Connection {
					ob.Subject = Subject[T]{}
					ob.Connection = nil
					ob.Disconnect = nil
					ob.ShareCount = 0
				}

				ob.Mu.Unlock()

				o.Emit(n)
			}
		})
	}

	c.AfterFunc(func() {
		ob.Mu.Lock()

		if connection == ob.Connection {
			ob.ShareCount--

			if ob.ShareCount == 0 {
				ob.Disconnect()

				ob.Subject = Subject[T]{}
				ob.Connection = nil
				ob.Disconnect = nil
			}
		}

		ob.Mu.Unlock()
	})
}
