package rx

import (
	"sync"
	"sync/atomic"
	"time"
)

// DebounceTime emits a value from the source [Observable] only after
// a particular time span has passed without another source emission.
func DebounceTime[T any](d time.Duration) Operator[T, T] {
	ob := Timer(d)
	return Debounce(func(T) Observable[time.Time] { return ob })
}

// Debounce emits a value from the source [Observable] only after a particular
// time span, determined by another [Observable], has passed without another
// source emission.
//
// It's like [DebounceTime], but the time span of emission silence is
// determined by a second [Observable].
func Debounce[T, U any](durationSelector func(v T) Observable[U]) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return debounceObservable[T, U]{source, durationSelector}.Subscribe
		},
	)
}

type debounceObservable[T, U any] struct {
	source           Observable[T]
	durationSelector func(T) Observable[U]
}

func (ob debounceObservable[T, U]) Subscribe(c Context, o Observer[T]) {
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
			sync.WaitGroup
			cancel CancelFunc
		}
	}

	x.context.Store(c.Context)

	startWorker := func(v T) {
		obs := ob.durationSelector(v)
		w, cancelw := c.WithCancel()

		x.context.Store(w.Context)
		x.worker.Add(1)
		x.worker.cancel = cancelw

		var noop bool

		obs.Subscribe(w, func(n Notification[U]) {
			if noop {
				return
			}

			defer x.worker.Done()

			noop = true
			cancelw()

			switch n.Kind {
			case KindNext:
				if x.latest.hasValue.Load() {
					x.latest.Lock()
					value := x.latest.value
					x.latest.hasValue.Store(false)
					x.latest.Unlock()
					Try1(o, Next(value), func() {
						if x.context.CompareAndSwap(w.Context, sentinel) {
							o.Stop(ErrOops)
						}
					})
				}

			case KindComplete:
				return

			case KindError, KindStop:
				if x.context.CompareAndSwap(w.Context, sentinel) {
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

	ob.source.Subscribe(c, func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
			x.latest.Lock()
			x.latest.value = n.Value
			x.latest.hasValue.Store(true)
			x.latest.Unlock()

			if x.context.Swap(c.Context) == sentinel {
				x.context.Store(sentinel)
				return
			}

			if x.worker.cancel != nil {
				x.worker.cancel()
				x.worker.Wait()
			}

			startWorker(n.Value)

		case KindComplete, KindError, KindStop:
			old := x.context.Swap(sentinel)

			cancel()
			x.worker.Wait()

			if old != sentinel {
				if n.Kind == KindComplete && x.latest.hasValue.Load() {
					Try1(o, Next(x.latest.value), func() { o.Stop(ErrOops) })
				}

				o.Emit(n)
			}
		}
	})
}
