package rx

import (
	"context"
	"runtime"
	"sync"
)

// An Observer is a consumer of notifications delivered by an [Observable].
type Observer[T any] func(n Notification[T])

// Noop gives you an Observer that does nothing.
func Noop[T any](Notification[T]) {}

// NewObserver creates an Observer from f.
func NewObserver[T any](f func(n Notification[T])) Observer[T] { return f }

// Next passes a value to sink.
func (sink Observer[T]) Next(v T) {
	sink(Next(v))
}

// Error passes an error to sink.
func (sink Observer[T]) Error(e error) {
	sink(Error[T](e))
}

// Complete passes a completion to sink.
func (sink Observer[T]) Complete() {
	sink(Complete[T]())
}

// Emit passes n to sink.
func (sink Observer[T]) Emit(n Notification[T]) {
	sink(n)
}

// ElementsOnly passes n to sink if n represents a value.
func (sink Observer[T]) ElementsOnly(n Notification[T]) {
	if n.Kind == KindNext {
		sink(n)
	}
}

// OnLastNotification creates an Observer that passes incoming emissions to
// sink, and when a notification of error or completion passes in, calls f
// just before passing it to sink.
func (sink Observer[T]) OnLastNotification(f func()) Observer[T] {
	return func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
			sink(n)
		case KindError, KindComplete:
			Try0(f, func() { sink.Error(ErrOops) })
			sink(n)
		}
	}
}

// Serialized creates an Observer that passes incoming emissions to sink
// in a mutually exclusive way.
func (sink Observer[T]) Serialized(c Context) Observer[T] {
	var x struct {
		Mu       sync.Mutex
		Emitting bool
		LastN    Notification[struct{}]
		Queue    []T
		Context  context.Context
		DoneChan <-chan struct{}
		Observer Observer[T]
	}

	x.Context = c.Context
	x.DoneChan = c.Done()
	x.Observer = sink

	return func(n Notification[T]) {
		x.Mu.Lock()

		if x.LastN.Kind != 0 {
			x.Mu.Unlock()
			return
		}

		switch n.Kind {
		case KindError:
			x.LastN = Error[struct{}](n.Error)
		case KindComplete:
			x.LastN = Complete[struct{}]()
		}

		if x.Emitting {
			if n.Kind == KindNext {
				x.Queue = append(x.Queue, n.Value)
			}

			x.Mu.Unlock()

			return
		}

		x.Emitting = true

		x.Mu.Unlock()

		throw := func(err error) {
			x.Mu.Lock()
			x.Emitting = false
			x.LastN = Error[struct{}](err)
			x.Queue = nil
			x.Mu.Unlock()
			x.Observer.Error(err)
		}

		oops := func() { throw(ErrOops) }

		sink := x.Observer

		switch n.Kind {
		case KindNext:
			select {
			default:
			case <-x.DoneChan:
				throw(x.Context.Err())
				return
			}

			Try1(sink, n, oops)

		case KindError, KindComplete:
			sink(n)
			return
		}

		for {
			x.Mu.Lock()

			if x.Queue == nil {
				lastn := x.LastN

				x.Emitting = false
				x.Mu.Unlock()

				switch lastn.Kind {
				case KindError:
					sink.Error(lastn.Error)
				case KindComplete:
					sink.Complete()
				}

				return
			}

			q := x.Queue
			x.Queue = nil

			x.Mu.Unlock()

			for _, v := range q {
				select {
				default:
				case <-x.DoneChan:
					throw(x.Context.Err())
					return
				}

				Try1(sink, Next(v), oops)
			}
		}
	}
}

// WithRuntimeFinalizer creates an Observer with a runtime finalizer set to
// run sink.Error(ErrFinalized) in a goroutine.
// sink must be safe for concurrent use.
func (sink Observer[T]) WithRuntimeFinalizer() Observer[T] {
	runtime.SetFinalizer(&sink, func(sink *Observer[T]) {
		go sink.Error(ErrFinalized)
	})

	return func(n Notification[T]) {
		switch n.Kind {
		case KindError, KindComplete:
			runtime.SetFinalizer(&sink, nil)
		}

		sink(n)
	}
}
