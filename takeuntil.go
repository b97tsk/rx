package rx

import (
	"context"
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

func (obs takeUntilObservable[T, U]) Subscribe(ctx context.Context, sink Observer[T]) {
	source, cancel := context.WithCancel(ctx)

	sink = sink.OnLastNotification(cancel)

	var x struct {
		Context atomic.Value
		Source  struct {
			sync.Mutex
			sync.WaitGroup
			Cancel context.CancelFunc
		}
		Worker struct {
			sync.WaitGroup
			Cancel context.CancelFunc
		}
	}

	x.Source.Cancel = cancel

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

			if (n.HasValue || n.HasError) && x.Context.CompareAndSwap(worker, sentinel) {
				x.Worker.Done()

				x.Source.Cancel()
				x.Source.Lock()
				x.Source.Wait()
				x.Source.Unlock()

				if n.HasValue {
					sink.Complete()
				} else {
					sink.Error(n.Error)
				}

				return
			}

			x.Worker.Done()
		})
	}

	x.Source.Lock()
	x.Source.Add(1)
	x.Source.Unlock()

	finish := func(n Notification[T]) {
		ctx := x.Context.Swap(source)

		x.Worker.Cancel()
		x.Worker.Wait()

		if x.Context.Swap(sentinel) != sentinel && ctx != sentinel {
			sink(n)
		}

		x.Source.Done()
	}

	if err := getErr(source); err != nil {
		finish(Error[T](err))
		return
	}

	obs.Source.Subscribe(source, func(n Notification[T]) {
		switch {
		case n.HasValue:
			sink(n)
		default:
			finish(n)
		}
	})
}
