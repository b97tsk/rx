package rx

import (
	"context"
	"sync"

	"github.com/b97tsk/rx/internal/queue"
)

// Unicast returns a Subject that only forwards every value it receives to
// the first subscriber.
// The first subscriber gets all future values.
// Subsequent subscribers will immediately receive an error notification of
// ErrUnicast.
// Values emitted to a Unicast before the first subscriber are lost.
func Unicast[T any]() Subject[T] {
	return UnicastBuffer[T](0)
}

// UnicastBufferAll returns a Subject that keeps track of every value
// it receives before the first subscriber.
// The first subscriber will then receive all tracked values as well as
// future values.
// Subsequent subscribers will immediately receive an error notification of
// ErrUnicast.
func UnicastBufferAll[T any]() Subject[T] {
	return UnicastBuffer[T](-1)
}

// UnicastBuffer returns a Subject that keeps track of a certain number of
// recent values it receive before the first subscriber.
// The first subscriber will then receive all tracked values as well as
// future values.
// Subsequent subscribers will immediately receive an error notification of
// ErrUnicast.
//
// If n < 0, UnicastBuffer keeps track of every value it receives before
// the first subscriber;
// if n == 0, UnicastBuffer doesn't keep track of any value it receives at all.
func UnicastBuffer[T any](n int) Subject[T] {
	u := &unicast[T]{Cap: n}
	return Subject[T]{
		Observable: NewObservable(u.subscribe),
		Observer:   NewObserver(u.emit).WithRuntimeFinalizer(),
	}
}

type unicast[T any] struct {
	Mu       sync.Mutex
	Cap      int
	Emitting bool
	LastN    Notification[struct{}]
	Buf      queue.Queue[T]
	Context  context.Context
	DoneChan <-chan struct{}
	Observer Observer[T]
}

func (u *unicast[T]) startEmitting(n Notification[T]) {
	u.Emitting = true
	u.Mu.Unlock()

	throw := func(err error) {
		u.Mu.Lock()
		u.Emitting = false
		u.LastN = Error[struct{}](err)
		u.Buf.Init()
		u.Mu.Unlock()
		u.Observer.Error(err)
	}

	oops := func() { throw(ErrOops) }

	sink := u.Observer

	switch n.Kind {
	case KindNext:
		select {
		default:
		case <-u.DoneChan:
			throw(u.Context.Err())
			return
		}

		Try1(sink, n, oops)

	case KindError, KindComplete:
		sink(n)
		return
	}

	for {
		u.Mu.Lock()

		if u.Buf.Len() == 0 {
			lastn := u.LastN

			u.Emitting = false
			u.Mu.Unlock()

			switch lastn.Kind {
			case KindError:
				sink.Error(lastn.Error)
			case KindComplete:
				sink.Complete()
			}

			return
		}

		var b queue.Queue[T]

		u.Buf, b = b, u.Buf

		u.Mu.Unlock()

		for i, j := 0, b.Len(); i < j; i++ {
			select {
			default:
			case <-u.DoneChan:
				throw(u.Context.Err())
				return
			}

			Try1(sink, Next(b.At(i)), oops)
		}
	}
}

func (u *unicast[T]) emit(n Notification[T]) {
	u.Mu.Lock()

	if u.LastN.Kind != 0 {
		u.Mu.Unlock()
		return
	}

	switch n.Kind {
	case KindError:
		u.LastN = Error[struct{}](n.Error)
	case KindComplete:
		u.LastN = Complete[struct{}]()
	}

	if u.Observer == nil {
		if n.Kind == KindNext && u.Cap != 0 {
			if u.Cap == u.Buf.Len() {
				u.Buf.Pop()
			}

			u.Buf.Push(n.Value)
		}

		u.Mu.Unlock()

		return
	}

	if u.Emitting {
		if n.Kind == KindNext {
			u.Buf.Push(n.Value)
		}

		u.Mu.Unlock()

		return
	}

	u.startEmitting(n)
}

func (u *unicast[T]) subscribe(c Context, sink Observer[T]) {
	u.Mu.Lock()

	if u.Observer != nil {
		defer u.Mu.Unlock()
		sink.Error(ErrUnicast)
		return
	}

	done := c.Done()
	if done != nil {
		stop := c.AfterFunc(func() { u.emit(Error[T](c.Err())) })
		sink = sink.OnLastNotification(func() { stop() })
	}

	u.Context = c.Context
	u.DoneChan = done
	u.Observer = sink

	if u.LastN.Kind == 0 && u.Buf.Len() == 0 {
		u.Mu.Unlock()
		return
	}

	u.startEmitting(Notification[T]{})
}
