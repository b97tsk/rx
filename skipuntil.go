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
	source, cancel := context.WithCancel(ctx)

	sink = sink.OnLastNotification(cancel)

	var x struct {
		Context atomic.Value
		Worker  struct {
			sync.WaitGroup
			Cancel context.CancelFunc
		}
	}

	{
		worker, cancel := context.WithCancel(source)
		x.Context.Store(worker)

		x.Worker.Cancel = cancel

		x.Worker.Add(1)

		var noop bool

		obs.Notifier.Subscribe(worker, func(n Notification[U]) {
			if noop {
				return
			}

			noop = true

			cancel()

			switch {
			case n.HasValue:
				x.Context.Store(source)

			case n.HasError:
				if x.Context.CompareAndSwap(worker, sentinel) {
					sink.Error(n.Error)
				}

			default:
				break
			}

			x.Worker.Done()
		})
	}

	finish := func(n Notification[T]) {
		ctx := x.Context.Swap(source)

		x.Worker.Cancel()
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
		switch {
		case n.HasValue:
			if x.Context.Load() == source {
				sink(n)
			}

		default:
			finish(n)
		}
	})
}
