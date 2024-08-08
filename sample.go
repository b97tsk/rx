package rx

import (
	"sync"
	"sync/atomic"
	"time"
)

// SampleTime emits the most recently emitted value from the source
// [Observalbe] within periodic time intervals.
func SampleTime[T any](d time.Duration) Operator[T, T] {
	return Sample[T](Ticker(d))
}

// Sample emits the most recently emitted value from the source [Observable]
// whenever notifier, another [Observable], emits a value.
func Sample[T, U any](notifier Observable[U]) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return sampleObservable[T, U]{source, notifier}.Subscribe
		},
	)
}

type sampleObservable[T, U any] struct {
	source   Observable[T]
	notifier Observable[U]
}

func (ob sampleObservable[T, U]) Subscribe(c Context, o Observer[T]) {
	c, cancel := c.WithCancel()
	o = o.DoOnTermination(cancel)

	var x struct {
		context atomic.Value
		latest  struct {
			sync.Mutex
			value    T
			hasValue atomic.Bool
		}
		worker struct {
			sync.Mutex
			sync.WaitGroup
		}
	}

	ob.source.Subscribe(c, func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
			x.latest.Lock()
			x.latest.value = n.Value
			x.latest.hasValue.Store(true)
			x.latest.Unlock()
		case KindComplete, KindError, KindStop:
			old := x.context.Swap(sentinel)

			cancel()

			x.worker.Lock()
			x.worker.Wait()
			x.worker.Unlock()

			if old != sentinel {
				o.Emit(n)
			}
		}
	})

	select {
	default:
	case <-c.Done():
		return
	}

	w, cancelw := c.WithCancel()

	x.worker.Lock()
	x.worker.Add(1)
	x.worker.Unlock()

	ob.notifier.Subscribe(w, func(n Notification[U]) {
		switch n.Kind {
		case KindNext:
			if x.latest.hasValue.Load() {
				x.latest.Lock()
				value := x.latest.value
				x.latest.hasValue.Store(false)
				x.latest.Unlock()
				o.Next(value)
			}

		case KindComplete:
			cancelw()
			x.worker.Done()

		case KindError, KindStop:
			defer x.worker.Done()

			cancelw()

			if x.context.Swap(sentinel) != sentinel {
				switch n.Kind {
				case KindError:
					o.Error(n.Error)
				case KindStop:
					o.Stop(n.Error)
				}
			}
		}
	})
}
