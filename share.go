package rx

import (
	"context"
	"sync"
)

// Share creates an [Observable] that multicasts (shares) the source
// [Observable]. When subscribed multiple times, it guarantees that only one
// subscription is made to the source at the same time. When all subscribers
// have unsubscribed it will unsubscribe from the source.
func Share[T any](c Context) ShareOperator[T] {
	return ShareOperator[T]{
		ts: shareConfig[T]{
			context:   c,
			connector: Multicast[T],
		},
	}
}

type shareConfig[T any] struct {
	context   Context
	connector func() Subject[T]
}

// ShareOperator is an [Operator] type for [Share].
type ShareOperator[T any] struct {
	ts shareConfig[T]
}

// WithConnector sets Connector option to a given value.
func (op ShareOperator[T]) WithConnector(connector func() Subject[T]) ShareOperator[T] {
	op.ts.connector = connector
	return op
}

// Apply implements the [Operator] interface.
func (op ShareOperator[T]) Apply(source Observable[T]) Observable[T] {
	ob := &shareObservable[T]{
		source:      source,
		shareConfig: op.ts,
	}
	return ob.Subscribe
}

type shareObservable[T any] struct {
	mu sync.Mutex

	source Observable[T]
	shareConfig[T]

	subject    Subject[T]
	connection context.Context
	disconnect CancelFunc
	shareCount int
}

func (ob *shareObservable[T]) Subscribe(c Context, o Observer[T]) {
	ob.mu.Lock()

	var unlocked bool

	defer func() {
		if !unlocked {
			ob.mu.Unlock()
		}
	}()

	if ob.subject.Observable == nil {
		ob.subject = Try01(ob.connector, func() { o.Stop(ErrOops) })
	}

	c, cancel := c.WithCancel()
	o = o.DoOnTermination(cancel)

	ob.subject.Subscribe(c, o)

	select {
	default:
	case <-c.Done():
		return
	}

	connection := ob.connection

	ob.shareCount++

	if connection == nil {
		w, cancelw := ob.context.WithCancel()

		connection = w.Context
		ob.connection = w.Context
		ob.disconnect = cancelw

		o := ob.subject.Observer

		ob.mu.Unlock()
		unlocked = true

		ob.source.Subscribe(w, func(n Notification[T]) {
			switch n.Kind {
			case KindNext:
				o.Emit(n)
			case KindComplete, KindError, KindStop:
				cancelw()

				ob.mu.Lock()

				if connection == ob.connection {
					ob.subject = Subject[T]{}
					ob.connection = nil
					ob.disconnect = nil
					ob.shareCount = 0
				}

				ob.mu.Unlock()

				o.Emit(n)
			}
		})
	}

	c.AfterFunc(func() {
		ob.mu.Lock()

		if connection == ob.connection {
			ob.shareCount--

			if ob.shareCount == 0 {
				ob.disconnect()

				ob.subject = Subject[T]{}
				ob.connection = nil
				ob.disconnect = nil
			}
		}

		ob.mu.Unlock()
	})
}
