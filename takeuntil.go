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
	source, cancelSource := context.WithCancel(ctx)

	sink = sink.OnLastNotification(cancelSource)

	var x struct {
		Context atomic.Value
		Source  struct {
			sync.Mutex
			sync.WaitGroup
		}
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
			case KindNext, KindError:
				if x.Context.CompareAndSwap(worker, sentinel) {
					cancelSource()

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

		cancelSource()

		if old != sentinel {
			sink(n)
		}

		x.Source.Done()
	}

	select {
	default:
	case <-source.Done():
		finish(Error[T](source.Err()))
		return
	}

	obs.Source.Subscribe(source, func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
			sink(n)
		case KindError, KindComplete:
			finish(n)
		}
	})
}
