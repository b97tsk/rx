package rx

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// Audit ignores source values for a duration determined by another Observable,
// then emits the most recent value from the source Observable, then repeats
// this process.
//
// It's like [AuditTime], but the silencing duration is determined by a second
// Observable.
func Audit[T, U any](durationSelector func(v T) Observable[U]) Operator[T, T] {
	if durationSelector == nil {
		panic("durationSelector == nil")
	}

	return audit(durationSelector)
}

// AuditTime ignores source values for a duration, then emits the most recent
// value from the source Observable, then repeats this process.
//
// When it sees a source value, it ignores that plus the next ones for a
// duration, and then it emits the most recent value from the source.
func AuditTime[T any](d time.Duration) Operator[T, T] {
	obsTimer := Timer(d)

	durationSelector := func(T) Observable[time.Time] { return obsTimer }

	return audit(durationSelector)
}

func audit[T, U any](durationSelector func(v T) Observable[U]) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return auditObservable[T, U]{source, durationSelector}.Subscribe
		},
	)
}

type auditObservable[T, U any] struct {
	Source           Observable[T]
	DurationSelector func(T) Observable[U]
}

func (obs auditObservable[T, U]) Subscribe(ctx context.Context, sink Observer[T]) {
	source, cancel := context.WithCancel(ctx)

	sink = sink.OnLastNotification(cancel)

	var x struct {
		Context  atomic.Value
		Complete atomic.Bool
		Latest   struct {
			sync.Mutex
			Value T
		}
		Worker struct {
			sync.WaitGroup
			Cancel context.CancelFunc
		}
	}

	x.Context.Store(source)

	startWorker := func(v T) {
		worker, cancel := context.WithCancel(source)
		x.Context.Store(worker)

		x.Worker.Cancel = cancel

		x.Worker.Add(1)

		var noop bool

		obs.DurationSelector(v).Subscribe(worker, func(n Notification[U]) {
			if noop {
				return
			}

			noop = true

			cancel()

			switch {
			case n.HasValue:
				x.Latest.Lock()
				value := x.Latest.Value
				x.Latest.Unlock()

				sink.Next(value)

				if x.Context.CompareAndSwap(worker, source) && x.Complete.Load() && x.Context.CompareAndSwap(source, sentinel) {
					sink.Complete()
				}

			case n.HasError:
				if x.Context.CompareAndSwap(worker, sentinel) {
					sink.Error(n.Error)
				}

			default:
				if x.Context.CompareAndSwap(worker, source) && x.Complete.Load() && x.Context.CompareAndSwap(source, sentinel) {
					sink.Complete()
				}
			}

			x.Worker.Done()
		})
	}

	obs.Source.Subscribe(source, func(n Notification[T]) {
		switch {
		case n.HasValue:
			x.Latest.Lock()
			x.Latest.Value = n.Value
			x.Latest.Unlock()

			if x.Context.Load() == source {
				startWorker(n.Value)
			}

		case n.HasError:
			ctx := x.Context.Swap(source)

			if x.Worker.Cancel != nil {
				x.Worker.Cancel()
				x.Worker.Wait()
			}

			if x.Context.Swap(sentinel) != sentinel && ctx != sentinel {
				sink(n)
			}

		default:
			x.Complete.Store(true)

			if x.Context.CompareAndSwap(source, sentinel) {
				sink(n)
			}
		}
	})
}
