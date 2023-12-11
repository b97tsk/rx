package rx

import (
	"context"
	"sync/atomic"
)

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

func (obs skipUntilObservable[T, U]) Subscribe(ctx context.Context, sink Observer[T]) {
	source, cancelSource := context.WithCancel(ctx)

	sink = sink.OnLastNotification(cancelSource)

	var x struct {
		Context atomic.Value
	}

	{
		worker, cancelWorker := context.WithCancel(source)

		x.Context.Store(worker)

		var noop bool

		obs.Notifier.Subscribe(worker, func(n Notification[U]) {
			if noop {
				return
			}

			noop = true

			cancelWorker()

			switch n.Kind {
			case KindNext:
				x.Context.CompareAndSwap(worker, source)

			case KindError:
				if x.Context.CompareAndSwap(worker, sentinel) {
					sink.Error(n.Error)
				}

			case KindComplete:
				break
			}
		})
	}

	finish := func(n Notification[T]) {
		old := x.Context.Swap(sentinel)

		cancelSource()

		if old != sentinel {
			sink(n)
		}
	}

	select {
	default:
	case <-source.Done():
		finish(Error[T](source.Err()))
		return
	}

	if x.Context.Load() == source {
		obs.Source.Subscribe(source, sink)
		return
	}

	obs.Source.Subscribe(source, func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
			if x.Context.Load() == source {
				sink(n)
			}

		case KindError, KindComplete:
			finish(n)
		}
	})
}
