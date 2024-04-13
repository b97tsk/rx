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

func (ob takeUntilObservable[T, U]) Subscribe(c Context, o Observer[T]) {
	c, cancel := c.WithCancel()
	o = o.DoOnTermination(cancel)

	var x struct {
		Context atomic.Value
		Source  struct {
			sync.Mutex
			sync.WaitGroup
		}
	}

	x.Context.Store(c.Context)

	{
		w, cancelw := c.WithCancel()

		var noop bool

		Try3(
			Observable[U].Subscribe,
			ob.Notifier,
			w,
			func(n Notification[U]) {
				if noop {
					return
				}

				noop = true
				cancelw()

				switch n.Kind {
				case KindNext, KindError:
					if x.Context.Swap(sentinel) != sentinel {
						cancel()

						x.Source.Lock()
						x.Source.Wait()
						x.Source.Unlock()

						switch n.Kind {
						case KindNext:
							o.Complete()
						case KindError:
							o.Error(n.Error)
						}
					}

				case KindComplete:
					return
				}
			},
			func() {
				if x.Context.Swap(sentinel) != sentinel {
					o.Error(ErrOops)
				}
			},
		)
	}

	x.Source.Lock()
	x.Source.Add(1)
	x.Source.Unlock()

	terminate := func(n Notification[T]) {
		defer x.Source.Done()

		old := x.Context.Swap(sentinel)

		cancel()

		if old != sentinel {
			o.Emit(n)
		}
	}

	select {
	default:
	case <-c.Done():
		terminate(Error[T](c.Cause()))
		return
	}

	ob.Source.Subscribe(c, func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
			if x.Context.Load() == c.Context {
				o.Emit(n)
			}
		case KindError, KindComplete:
			terminate(n)
		}
	})
}
