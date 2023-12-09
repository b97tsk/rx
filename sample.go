package rx

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// SampleTime emits the most recently emitted value from the source Observalbe
// within periodic time intervals.
func SampleTime[T any](d time.Duration) Operator[T, T] {
	return Sample[T](Ticker(d))
}

// Sample emits the most recently emitted value from the source Observable
// whenever notifier, another Observable, emits a value.
func Sample[T, U any](notifier Observable[U]) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return sampleObservable[T, U]{source, notifier}.Subscribe
		},
	)
}

type sampleObservable[T, U any] struct {
	Source   Observable[T]
	Notifier Observable[U]
}

func (obs sampleObservable[T, U]) Subscribe(ctx context.Context, sink Observer[T]) {
	source, cancelSource := context.WithCancel(ctx)

	sink = sink.OnLastNotification(cancelSource)

	var x struct {
		Context atomic.Value
		Latest  struct {
			sync.Mutex
			Value    T
			HasValue atomic.Bool
		}
		Worker struct {
			sync.WaitGroup
		}
	}

	{
		worker, cancelWorker := context.WithCancel(source)
		x.Context.Store(worker)

		x.Worker.Add(1)

		obs.Notifier.Subscribe(worker, func(n Notification[U]) {
			switch n.Kind {
			case KindNext:
				if x.Latest.HasValue.Load() {
					x.Latest.Lock()
					value := x.Latest.Value
					x.Latest.HasValue.Store(false)
					x.Latest.Unlock()
					sink.Next(value)
				}

				return

			case KindError:
				if x.Context.CompareAndSwap(worker, sentinel) {
					sink.Error(n.Error)
				}

			case KindComplete:
				break
			}

			cancelWorker()
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

	obs.Source.Subscribe(source, func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
			x.Latest.Lock()
			x.Latest.Value = n.Value
			x.Latest.HasValue.Store(true)
			x.Latest.Unlock()

		case KindError, KindComplete:
			finish(n)
		}
	})
}
