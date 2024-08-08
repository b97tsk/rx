package rx

import (
	"sync"
	"sync/atomic"
)

// TakeUntil mirrors the source [Observable] until a second [Observable] emits
// a value.
func TakeUntil[T, U any](notifier Observable[U]) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return takeUntilObservable[T, U]{source, notifier}.Subscribe
		},
	)
}

type takeUntilObservable[T, U any] struct {
	source   Observable[T]
	notifier Observable[U]
}

func (ob takeUntilObservable[T, U]) Subscribe(c Context, o Observer[T]) {
	c, cancel := c.WithCancel()
	o = o.DoOnTermination(cancel)

	var x struct {
		context atomic.Value
		source  struct {
			sync.Mutex
			sync.WaitGroup
		}
	}

	x.context.Store(c.Context)

	{
		w, cancelw := c.WithCancel()

		var noop bool

		Try3(
			Observable[U].Subscribe,
			ob.notifier,
			w,
			func(n Notification[U]) {
				if noop {
					return
				}

				noop = true
				cancelw()

				switch n.Kind {
				case KindNext, KindError, KindStop:
					if x.context.Swap(sentinel) != sentinel {
						cancel()

						x.source.Lock()
						x.source.Wait()
						x.source.Unlock()

						switch n.Kind {
						case KindNext:
							o.Complete()
						case KindError:
							o.Error(n.Error)
						case KindStop:
							o.Stop(n.Error)
						}
					}

				case KindComplete:
					return
				}
			},
			func() {
				if x.context.Swap(sentinel) != sentinel {
					o.Stop(ErrOops)
				}
			},
		)
	}

	x.source.Lock()
	x.source.Add(1)
	x.source.Unlock()

	terminate := func(n Notification[T]) {
		defer x.source.Done()

		old := x.context.Swap(sentinel)

		cancel()

		if old != sentinel {
			o.Emit(n)
		}
	}

	select {
	default:
	case <-c.Done():
		terminate(Stop[T](c.Cause()))
		return
	}

	ob.source.Subscribe(c, func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
			if x.context.Load() == c.Context {
				o.Emit(n)
			}
		case KindComplete, KindError, KindStop:
			terminate(n)
		}
	})
}
