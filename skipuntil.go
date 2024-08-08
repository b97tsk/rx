package rx

import "sync/atomic"

// SkipUntil skips values emitted by the source [Observable] until
// a second [Observable] emits an value.
func SkipUntil[T, U any](notifier Observable[U]) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return skipUntilObservable[T, U]{source, notifier}.Subscribe
		},
	)
}

type skipUntilObservable[T, U any] struct {
	source   Observable[T]
	notifier Observable[U]
}

func (ob skipUntilObservable[T, U]) Subscribe(c Context, o Observer[T]) {
	c, cancel := c.WithCancel()
	o = o.DoOnTermination(cancel)

	var x struct {
		context atomic.Value
	}

	{
		w, cancelw := c.WithCancel()

		x.context.Store(w.Context)

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
				case KindNext:
					x.context.CompareAndSwap(w.Context, c.Context)

				case KindComplete:
					return

				case KindError, KindStop:
					if x.context.CompareAndSwap(w.Context, sentinel) {
						switch n.Kind {
						case KindError:
							o.Error(n.Error)
						case KindStop:
							o.Stop(n.Error)
						}
					}
				}
			},
			func() {
				if x.context.Swap(sentinel) != sentinel {
					o.Stop(ErrOops)
				}
			},
		)
	}

	terminate := func(n Notification[T]) {
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

	if x.context.Load() == c.Context {
		ob.source.Subscribe(c, o)
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
