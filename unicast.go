package rx

import (
	"context"
	"sync"

	"github.com/b97tsk/rx/internal/queue"
)

// Unicast returns a [Subject] that forwards every value it receives to
// the first subscriber.
// The first subscriber gets all future values.
// Subsequent subscribers will immediately receive a [Stop] notification of
// [ErrUnicast].
// Values emitted to a Unicast before the first subscriber are lost.
func Unicast[T any]() Subject[T] {
	return UnicastBuffer[T](0)
}

// UnicastBufferAll returns a [Subject] that keeps track of every value
// it receives before the first subscriber.
// The first subscriber will then receive all tracked values as well as
// future values.
// Subsequent subscribers will immediately receive a [Stop] notification of
// [ErrUnicast].
func UnicastBufferAll[T any]() Subject[T] {
	return UnicastBuffer[T](-1)
}

// UnicastBuffer returns a [Subject] that keeps track of a certain number of
// recent values it receive before the first subscriber.
// The first subscriber will then receive all tracked values as well as
// future values.
// Subsequent subscribers will immediately receive a [Stop] notification of
// [ErrUnicast].
//
// If n < 0, UnicastBuffer keeps track of every value it receives before
// the first subscriber;
// if n == 0, UnicastBuffer doesn't keep track of any value it receives at all.
func UnicastBuffer[T any](n int) Subject[T] {
	u := &unicast[T]{cap: n}
	return Subject[T]{
		Observable: NewObservable(u.Subscribe),
		Observer:   WithRuntimeFinalizer(u.Emit),
	}
}

type unicast[T any] struct {
	mu       sync.Mutex
	co       sync.Cond
	cap      int
	waiters  int
	emitting bool
	lastn    Notification[struct{}]
	buf, alt queue.Queue[T]
	context  context.Context
	done     <-chan struct{}
	observer Observer[T]
}

const bufLimit = 32

func (u *unicast[T]) Emit(n Notification[T]) {
	u.mu.Lock()

	if u.lastn.Kind != 0 {
		u.mu.Unlock()
		return
	}

	switch n.Kind {
	case KindNext:
	case KindComplete:
		u.lastn = Complete[struct{}]()
	case KindError:
		u.lastn = Error[struct{}](n.Error)
	case KindStop:
		u.lastn = Stop[struct{}](n.Error)
	}

	if u.observer == nil {
		if n.Kind == KindNext && u.cap != 0 {
			if u.cap == u.buf.Len() {
				u.buf.Pop()
			}

			u.buf.Push(n.Value)
		}

		u.mu.Unlock()

		return
	}

	if !u.emitting {
		u.startEmitting(n)
		return
	}

	if n.Kind == KindNext {
	Again:
		for u.buf.Len() >= max(u.buf.Cap(), bufLimit) {
			if n.Kind == 0 && u.waiters != 0 {
				break
			}

			if !u.emitting {
				u.startEmitting(n)
				return
			}

			if u.co.L == nil {
				u.co.L = &u.mu
			}

			u.waiters++
			u.co.Wait()
			u.waiters--
		}

		if n.Kind == KindNext && u.lastn.Kind == 0 {
			u.buf.Push(n.Value)
			if u.waiters == 0 {
				n = Notification[T]{}
				goto Again
			}
		}
	}

	u.mu.Unlock()
}

func (u *unicast[T]) Subscribe(c Context, o Observer[T]) {
	u.mu.Lock()

	if u.observer != nil {
		u.mu.Unlock()
		o.Stop(ErrUnicast)
		return
	}

	done := c.Done()
	if done != nil {
		stop := c.AfterFunc(func() { u.Emit(Stop[T](c.Cause())) })
		o = o.DoOnTermination(func() { stop() })
	}

	u.context = c.Context
	u.done = done
	u.observer = o

	if u.lastn.Kind == 0 && u.buf.Len() == 0 {
		u.mu.Unlock()
		return
	}

	u.startEmitting(Notification[T]{})
}

func (u *unicast[T]) startEmitting(n Notification[T]) {
	stop := func(err error) {
		u.mu.Lock()
		u.emitting = false
		u.lastn = Stop[struct{}](err)
		u.buf.Init()
		u.alt.Init()
		u.mu.Unlock()
		u.co.Broadcast()
		u.observer.Stop(err)
	}

	oops := func() { stop(ErrOops) }

	o := u.observer

	u.emitting = true

	for {
		var buf queue.Queue[T]

		buf, u.buf, u.alt = u.buf, u.alt, buf

		u.mu.Unlock()
		u.co.Broadcast()

		for i := range buf.Len() {
			select {
			default:
			case <-u.done:
				stop(context.Cause(u.context))
				return
			}

			Try1(o, Next(buf.At(i)), oops)
		}

		switch n.Kind {
		case 0:
		case KindNext:
			select {
			default:
			case <-u.done:
				stop(context.Cause(u.context))
				return
			}

			Try1(o, n, oops)

			n = Notification[T]{}
		case KindComplete, KindError, KindStop:
			o.Emit(n)
			return
		}

		u.mu.Lock()

		if n := buf.Cap(); n != 0 && n <= bufLimit {
			buf.Clear()
			if u.buf.Cap() == 0 {
				u.buf = buf
			} else {
				u.alt = buf
			}
		}

		if u.buf.Len() == 0 {
			lastn := u.lastn

			u.emitting = false
			u.mu.Unlock()

			switch lastn.Kind {
			case KindComplete:
				o.Complete()
			case KindError:
				o.Error(lastn.Error)
			case KindStop:
				o.Stop(lastn.Error)
			}

			return
		}

		if u.waiters != 0 {
			u.emitting = false
			u.mu.Unlock()
			u.co.Broadcast()
			return
		}
	}
}
