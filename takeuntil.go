package rx

import (
	"sync"
	"sync/atomic"
)

// TakeUntil mirrors the source Observable until a second Observable emits
// a value.
func TakeUntil[T, U any](notifier Observable[U]) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return takeUntilObservable[T, U]{source, notifier}.Subscribe
		},
	)
}

type takeUntilObservable[T, U any] struct {
	Source   Observable[T]
	Notifier Observable[U]
}

func (obs takeUntilObservable[T, U]) Subscribe(c Context, sink Observer[T]) {
	c, cancel := c.WithCancel()
	sink = sink.OnLastNotification(cancel)

	var x struct {
		Context atomic.Value
		Source  struct {
			sync.Mutex
			sync.WaitGroup
		}
	}

	{
		w, cancelw := c.WithCancel()

		x.Context.Store(w.Context)

		var noop bool

		obs.Notifier.Subscribe(w, func(n Notification[U]) {
			if noop {
				return
			}

			noop = true
			cancelw()

			switch n.Kind {
			case KindNext, KindError:
				if x.Context.CompareAndSwap(w.Context, sentinel) {
					cancel()

					x.Source.Lock()
					x.Source.Wait()
					x.Source.Unlock()

					switch n.Kind {
					case KindNext:
						sink.Complete()
					case KindError:
						sink.Error(n.Error)
					}

					return
				}
			}
		})
	}

	x.Source.Lock()
	x.Source.Add(1)
	x.Source.Unlock()

	finish := func(n Notification[T]) {
		old := x.Context.Swap(sentinel)

		cancel()

		if old != sentinel {
			sink(n)
		}

		x.Source.Done()
	}

	select {
	default:
	case <-c.Done():
		finish(Error[T](c.Err()))
		return
	}

	obs.Source.Subscribe(c, func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
			sink(n)
		case KindError, KindComplete:
			finish(n)
		}
	})
}
