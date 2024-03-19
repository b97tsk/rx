package rx

import (
	"context"
	"sync"
)

// Serialize returns an Observer that passes incoming emissions to sink
// in a mutually exclusive way.
// Serialize also returns a copy of c that will be cancelled when sink is
// about to receive a notification of error or completion.
func Serialize[T any](c Context, sink Observer[T]) (Context, Observer[T]) {
	c, cancel := c.WithCancel()

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
	x.Observer = sink.OnLastNotification(cancel)

	return c, func(n Notification[T]) {
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
