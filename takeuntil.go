package rx

import (
	"context"
	"errors"
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
	c, cancel := c.WithCancelCause()
	o = o.DoOnTermination(func() { cancel(nil) })

	var x struct {
		Context atomic.Value
	}

	x.Context.Store(c.Context)

	complete := errors.Join(context.Canceled)

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
				case KindNext:
					if x.Context.CompareAndSwap(c.Context, w.Context) {
						cancel(complete)
					}
				case KindError:
					if x.Context.CompareAndSwap(c.Context, w.Context) {
						cancel(n.Error)
					}
				case KindComplete:
					return
				}
			},
			func() { o.Error(ErrOops) },
		)
	}

	terminate := func(n Notification[T]) {
		if n.Kind == KindError && errors.Is(n.Error, complete) {
			o.Complete()
			return
		}

		o.Emit(n)
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
