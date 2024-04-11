package rx

import (
	"sync"
	"sync/atomic"
	"time"
)

// ThrottleTime emits a value from the source Observable, then ignores
// subsequent source values for a duration, then repeats this process until
// the source completes.
//
// ThrottleTime lets a value pass, then ignores source values
// for the next duration time.
func ThrottleTime[T any](d time.Duration) ThrottleOperator[T, time.Time] {
	ob := Timer(d)
	return Throttle(func(T) Observable[time.Time] { return ob })
}

// Throttle emits a value from the source Observable, then ignores
// subsequent source values for a duration determined by another Observable,
// then repeats this process until the source completes.
//
// It's like [ThrottleTime], but the silencing duration is determined
// by a second Observable.
func Throttle[T, U any](durationSelector func(v T) Observable[U]) ThrottleOperator[T, U] {
	return ThrottleOperator[T, U]{
		ts: throttleConfig[T, U]{
			DurationSelector: durationSelector,
			Leading:          true,
			Trailing:         false,
		},
	}
}

type throttleConfig[T, U any] struct {
	DurationSelector func(T) Observable[U]
	Leading          bool
	Trailing         bool
}

// ThrottleOperator is an [Operator] type for [Throttle].
type ThrottleOperator[T, U any] struct {
	ts throttleConfig[T, U]
}

// WithLeading sets Leading option to a given value.
func (op ThrottleOperator[T, U]) WithLeading(v bool) ThrottleOperator[T, U] {
	op.ts.Leading = v
	return op
}

// WithTrailing sets Trailing option to a given value.
func (op ThrottleOperator[T, U]) WithTrailing(v bool) ThrottleOperator[T, U] {
	op.ts.Trailing = v
	return op
}

// Apply implements the Operator interface.
func (op ThrottleOperator[T, U]) Apply(source Observable[T]) Observable[T] {
	return throttleObservable[T, U]{source, op.ts}.Subscribe
}

type throttleObservable[T, U any] struct {
	Source Observable[T]
	throttleConfig[T, U]
}

func (ob throttleObservable[T, U]) Subscribe(c Context, o Observer[T]) {
	c, cancel := c.WithCancel()
	o = o.DoOnTermination(cancel)

	var x struct {
		Context  atomic.Value
		Complete atomic.Bool
		Trailing struct {
			sync.Mutex
			Value    T
			HasValue atomic.Bool
		}
		Worker struct {
			sync.WaitGroup
		}
	}

	x.Context.Store(c.Context)

	var doThrottle func(T)

	doThrottle = func(v T) {
		obs := ob.DurationSelector(v)
		w, cancelw := c.WithCancel()

		x.Context.Store(w.Context)
		x.Worker.Add(1)

		var noop bool

		obs.Subscribe(w, func(n Notification[U]) {
			if noop {
				return
			}

			defer x.Worker.Done()

			noop = true
			cancelw()

			switch n.Kind {
			case KindNext:
				if ob.Trailing && x.Trailing.HasValue.Load() {
					x.Trailing.Lock()
					value := x.Trailing.Value
					x.Trailing.HasValue.Store(false)
					x.Trailing.Unlock()

					oops := func() {
						if x.Context.Swap(sentinel) != sentinel {
							o.Error(ErrOops)
						}
					}

					Try1(o, Next(value), oops)

					if !x.Complete.Load() {
						Try1(doThrottle, value, oops)
					}
				}

				if x.Context.CompareAndSwap(w.Context, c.Context) && x.Complete.Load() && x.Context.CompareAndSwap(c.Context, sentinel) {
					o.Complete()
				}

			case KindError:
				if x.Context.Swap(sentinel) != sentinel {
					o.Error(n.Error)
				}

			case KindComplete:
				if x.Context.CompareAndSwap(w.Context, c.Context) && x.Complete.Load() && x.Context.CompareAndSwap(c.Context, sentinel) {
					o.Complete()
				}
			}
		})
	}

	ob.Source.Subscribe(c, func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
			x.Trailing.Lock()
			x.Trailing.Value = n.Value
			x.Trailing.HasValue.Store(true)
			x.Trailing.Unlock()

			if x.Context.Load() == c.Context {
				if ob.Leading {
					x.Trailing.HasValue.Store(false)
					o.Emit(n)
				}

				doThrottle(n.Value)
			}

		case KindError:
			old := x.Context.Swap(sentinel)

			cancel()
			x.Worker.Wait()

			if old != sentinel {
				o.Emit(n)
			}

		case KindComplete:
			x.Complete.Store(true)

			if x.Context.CompareAndSwap(c.Context, sentinel) {
				o.Emit(n)
			}
		}
	})
}
