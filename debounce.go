package rx

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// Debounce emits a value from the source Observable only after a particular
// time span, determined by another Observable, has passed without another
// source emission.
//
// It's like [DebounceTime], but the time span of emission silence is
// determined by a second Observable.
func Debounce[T, U any](durationSelector func(v T) Observable[U]) Operator[T, T] {
	if durationSelector == nil {
		panic("durationSelector == nil")
	}

	return debounce(durationSelector)
}

// DebounceTime emits a value from the source Observable only after
// a particular time span has passed without another source emission.
func DebounceTime[T any](d time.Duration) Operator[T, T] {
	obsTimer := Timer(d)

	durationSelector := func(T) Observable[time.Time] { return obsTimer }

	return debounce(durationSelector)
}

func debounce[T, U any](durationSelector func(v T) Observable[U]) Operator[T, T] {
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

func (obs debounceObservable[T, U]) Subscribe(ctx context.Context, sink Observer[T]) {
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
			Cancel context.CancelFunc
		}
	}

	x.Context.Store(source)

	startWorker := func(v T) {
		worker, cancelWorker := context.WithCancel(source)

		x.Context.Store(worker)

		x.Worker.Add(1)

		x.Worker.Cancel = cancelWorker

		var noop bool

		obs.DurationSelector(v).Subscribe(worker, func(n Notification[U]) {
			if noop {
				return
			}

			noop = true

			cancelWorker()

			switch n.Kind {
			case KindNext:
				if x.Latest.HasValue.Load() {
					x.Latest.Lock()
					value := x.Latest.Value
					x.Latest.HasValue.Store(false)
					x.Latest.Unlock()
					sink.Next(value)
				}

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

	obs.Source.Subscribe(source, func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
			x.Latest.Lock()
			x.Latest.Value = n.Value
			x.Latest.HasValue.Store(true)
			x.Latest.Unlock()

			if x.Context.Swap(source) == sentinel {
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

			cancelSource()
			x.Worker.Wait()

			if old != sentinel {
				if n.Kind == KindComplete && x.Latest.HasValue.Load() {
					sink.Next(x.Latest.Value)
				}

				sink(n)
			}
		}
	})
}
