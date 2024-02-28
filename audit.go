package rx

import (
	"sync"
	"sync/atomic"
	"time"
)

// AuditTime ignores source values for a duration, then emits the most recent
// value from the source Observable, then repeats this process.
//
// When it sees a source value, it ignores that plus the next ones for a
// duration, and then it emits the most recent value from the source.
func AuditTime[T any](d time.Duration) Operator[T, T] {
	obsTimer := Timer(d)
	durationSelector := func(T) Observable[time.Time] { return obsTimer }
	return Audit(durationSelector)
}

// Audit ignores source values for a duration determined by another Observable,
// then emits the most recent value from the source Observable, then repeats
// this process.
//
// It's like [AuditTime], but the silencing duration is determined by a second
// Observable.
func Audit[T, U any](durationSelector func(v T) Observable[U]) Operator[T, T] {
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

func (obs auditObservable[T, U]) Subscribe(c Context, sink Observer[T]) {
	c, cancel := c.WithCancel()
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
		}
	}

	x.Context.Store(c.Context)

	startWorker := func(v T) {
		obs1 := obs.DurationSelector(v)
		w, cancelw := c.WithCancel()

		x.Context.Store(w.Context)
		x.Worker.Add(1)

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
				x.Latest.Lock()
				value := x.Latest.Value
				x.Latest.Unlock()

				Try1(sink, Next(value), func() {
					if x.Context.Swap(sentinel) != sentinel {
						sink.Error(ErrOops)
					}
				})

				if x.Context.CompareAndSwap(w.Context, c.Context) && x.Complete.Load() && x.Context.CompareAndSwap(c.Context, sentinel) {
					sink.Complete()
				}

			case KindError:
				if x.Context.Swap(sentinel) != sentinel {
					sink.Error(n.Error)
				}

			case KindComplete:
				if x.Context.CompareAndSwap(w.Context, c.Context) && x.Complete.Load() && x.Context.CompareAndSwap(c.Context, sentinel) {
					sink.Complete()
				}
			}
		})
	}

	obs.Source.Subscribe(c, func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
			x.Latest.Lock()
			x.Latest.Value = n.Value
			x.Latest.Unlock()

			if x.Context.Load() == c.Context {
				startWorker(n.Value)
			}

		case KindError:
			old := x.Context.Swap(sentinel)

			cancel()
			x.Worker.Wait()

			if old != sentinel {
				sink(n)
			}

		case KindComplete:
			x.Complete.Store(true)

			if x.Context.CompareAndSwap(c.Context, sentinel) {
				sink(n)
			}
		}
	})
}
