package rx

import (
	"sync"
	"sync/atomic"
	"time"
)

// DebounceTime emits a value from the source Observable only after
// a particular time span has passed without another source emission.
func DebounceTime[T any](d time.Duration) Operator[T, T] {
	obsTimer := Timer(d)
	durationSelector := func(T) Observable[time.Time] { return obsTimer }
	return Debounce(durationSelector)
}

// Debounce emits a value from the source Observable only after a particular
// time span, determined by another Observable, has passed without another
// source emission.
//
// It's like [DebounceTime], but the time span of emission silence is
// determined by a second Observable.
func Debounce[T, U any](durationSelector func(v T) Observable[U]) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return debounceObservable[T, U]{source, durationSelector}.Subscribe
		},
	)
}

type debounceObservable[T, U any] struct {
	Source           Observable[T]
	DurationSelector func(T) Observable[U]
}

func (obs debounceObservable[T, U]) Subscribe(c Context, sink Observer[T]) {
	c, cancel := c.WithCancel()
	sink = sink.OnTermination(cancel)

	var x struct {
		Context atomic.Value
		Latest  struct {
			sync.Mutex
			Value    T
			HasValue atomic.Bool
		}
		Worker struct {
			sync.WaitGroup
			Cancel CancelFunc
		}
	}

	x.Context.Store(c.Context)

	startWorker := func(v T) {
		obs1 := obs.DurationSelector(v)
		w, cancelw := c.WithCancel()

		x.Context.Store(w.Context)
		x.Worker.Add(1)
		x.Worker.Cancel = cancelw

		var noop bool

		obs1.Subscribe(w, func(n Notification[U]) {
			if noop {
				return
			}

			defer x.Worker.Done()

			noop = true
			cancelw()

			switch n.Kind {
			case KindNext:
				if x.Latest.HasValue.Load() {
					x.Latest.Lock()
					value := x.Latest.Value
					x.Latest.HasValue.Store(false)
					x.Latest.Unlock()
					Try1(sink, Next(value), func() {
						if x.Context.CompareAndSwap(w.Context, sentinel) {
							sink.Error(ErrOops)
						}
					})
				}

			case KindError:
				if x.Context.CompareAndSwap(w.Context, sentinel) {
					sink.Error(n.Error)
				}

			case KindComplete:
				return
			}
		})
	}

	obs.Source.Subscribe(c, func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
			x.Latest.Lock()
			x.Latest.Value = n.Value
			x.Latest.HasValue.Store(true)
			x.Latest.Unlock()

			if x.Context.Swap(c.Context) == sentinel {
				x.Context.Store(sentinel)
				return
			}

			if x.Worker.Cancel != nil {
				x.Worker.Cancel()
				x.Worker.Wait()
			}

			startWorker(n.Value)

		case KindError, KindComplete:
			old := x.Context.Swap(sentinel)

			cancel()
			x.Worker.Wait()

			if old != sentinel {
				if n.Kind == KindComplete && x.Latest.HasValue.Load() {
					Try1(sink, Next(x.Latest.Value), func() { sink.Error(ErrOops) })
				}

				sink(n)
			}
		}
	})
}
