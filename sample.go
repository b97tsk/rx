package rx

import (
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

func (obs sampleObservable[T, U]) Subscribe(c Context, sink Observer[T]) {
	c, cancel := c.WithCancel()
	sink = sink.OnLastNotification(cancel)

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
		w, cancelw := c.WithCancel()

		x.Context.Store(w.Context)
		x.Worker.Add(1)

		obs.Notifier.Subscribe(w, func(n Notification[U]) {
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
				if x.Context.CompareAndSwap(w.Context, sentinel) {
					sink.Error(n.Error)
				}

			case KindComplete:
				break
			}

			cancelw()
			x.Worker.Done()
		})
	}

	finish := func(n Notification[T]) {
		old := x.Context.Swap(sentinel)

		cancel()
		x.Worker.Wait()

		if old != sentinel {
			sink(n)
		}
	}

	select {
	default:
	case <-c.Done():
		finish(Error[T](c.Err()))
		return
	}

	obs.Source.Subscribe(c, func(n Notification[T]) {
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
