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
	ob := Timer(d)
	return Audit(func(T) Observable[time.Time] { return ob })
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
	source           Observable[T]
	durationSelector func(T) Observable[U]
}

func (ob auditObservable[T, U]) Subscribe(c Context, o Observer[T]) {
	c, cancel := c.WithCancel()
	o = o.DoOnTermination(cancel)

	var x struct {
		context  atomic.Value
		complete atomic.Bool
		latest   struct {
			sync.Mutex
			value T
		}
		worker struct {
			sync.WaitGroup
		}
	}

	x.context.Store(c.Context)

	startWorker := func(v T) {
		obs := ob.durationSelector(v)
		w, cancelw := c.WithCancel()

		x.context.Store(w.Context)
		x.worker.Add(1)

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
				x.latest.Lock()
				value := x.latest.value
				x.latest.Unlock()

				Try1(o, Next(value), func() {
					if x.context.Swap(sentinel) != sentinel {
						o.Error(ErrOops)
					}
				})

				if x.context.CompareAndSwap(w.Context, c.Context) && x.complete.Load() && x.context.CompareAndSwap(c.Context, sentinel) {
					o.Complete()
				}

			case KindError:
				if x.context.Swap(sentinel) != sentinel {
					o.Error(n.Error)
				}

			case KindComplete:
				if x.context.CompareAndSwap(w.Context, c.Context) && x.complete.Load() && x.context.CompareAndSwap(c.Context, sentinel) {
					o.Complete()
				}
			}
		})
	}

	ob.source.Subscribe(c, func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
			x.latest.Lock()
			x.latest.value = n.Value
			x.latest.Unlock()

			if x.context.Load() == c.Context {
				startWorker(n.Value)
			}

		case KindError:
			old := x.context.Swap(sentinel)

			cancel()
			x.worker.Wait()

			if old != sentinel {
				o.Emit(n)
			}

		case KindComplete:
			x.complete.Store(true)

			if x.context.CompareAndSwap(c.Context, sentinel) {
				o.Emit(n)
			}
		}
	})
}
