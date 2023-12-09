package rx

import (
	"context"
	"sync"
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
		Worker  struct {
			sync.WaitGroup
		}
	}

	{
		worker, cancelWorker := context.WithCancel(source)
		x.Context.Store(worker)

		x.Worker.Add(1)

		var noop bool

		obs.Notifier.Subscribe(worker, func(n Notification[U]) {
			if noop {
				return
			}

			noop = true

			cancelWorker()

			switch n.Kind {
			case KindNext:
				x.Context.Store(source)

			case KindError:
				if x.Context.CompareAndSwap(worker, sentinel) {
					sink.Error(n.Error)
				}

			case KindComplete:
				break
			}

			x.Worker.Done()
		})
	}

	finish := func(n Notification[T]) {
		ctx := x.Context.Swap(source)

		cancelSource()
		x.Worker.Wait()

		if x.Context.Swap(sentinel) != sentinel && ctx != sentinel {
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
