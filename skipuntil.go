package rx

import "sync/atomic"

// SkipUntil skips values emitted by the source Observable
// until a second Observable emits an value.
func SkipUntil[T, U any](notifier Observable[U]) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return skipUntilObservable[T, U]{source, notifier}.Subscribe
		},
	)
}

type skipUntilObservable[T, U any] struct {
	Source   Observable[T]
	Notifier Observable[U]
}

func (obs skipUntilObservable[T, U]) Subscribe(c Context, sink Observer[T]) {
	c, cancel := c.WithCancel()
	sink = sink.DoOnTermination(cancel)

	var x struct {
		Context atomic.Value
	}

	{
		w, cancelw := c.WithCancel()

		x.Context.Store(w.Context)

		var noop bool

		Try3(
			Observable[U].Subscribe,
			obs.Notifier,
			w,
			func(n Notification[U]) {
				if noop {
					return
				}

				noop = true
				cancelw()

				switch n.Kind {
				case KindNext:
					x.Context.CompareAndSwap(w.Context, c.Context)

				case KindError:
					if x.Context.CompareAndSwap(w.Context, sentinel) {
						sink.Error(n.Error)
					}

				case KindComplete:
					return
				}
			},
			func() {
				if x.Context.Swap(sentinel) != sentinel {
					sink.Error(ErrOops)
				}
			},
		)
	}

	terminate := func(n Notification[T]) {
		old := x.Context.Swap(sentinel)

		cancel()

		if old != sentinel {
			sink(n)
		}
	}

	select {
	default:
	case <-c.Done():
		terminate(Error[T](c.Err()))
		return
	}

	if x.Context.Load() == c.Context {
		obs.Source.Subscribe(c, sink)
		return
	}

	obs.Source.Subscribe(c, func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
			if x.Context.Load() == c.Context {
				sink(n)
			}
		case KindError, KindComplete:
			terminate(n)
		}
	})
}
