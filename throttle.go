package rx

import (
	"sync"
	"sync/atomic"
	"time"
)

// ThrottleTime emits a value from the source [Observable], then ignores
// subsequent source values for a duration, then repeats this process until
// the source completes.
//
// ThrottleTime lets a value pass, then ignores source values
// for the next duration time.
func ThrottleTime[T any](d time.Duration) ThrottleOperator[T, time.Time] {
	ob := Timer(d)
	return Throttle(func(T) Observable[time.Time] { return ob })
}

// Throttle emits a value from the source [Observable], then ignores
// subsequent source values for a duration determined by another [Observable],
// then repeats this process until the source completes.
//
// It's like [ThrottleTime], but the silencing duration is determined
// by a second [Observable].
func Throttle[T, U any](durationSelector func(v T) Observable[U]) ThrottleOperator[T, U] {
	return ThrottleOperator[T, U]{
		ts: throttleConfig[T, U]{
			durationSelector: durationSelector,
			leading:          true,
			trailing:         false,
		},
	}
}

type throttleConfig[T, U any] struct {
	durationSelector func(T) Observable[U]
	leading          bool
	trailing         bool
}

// ThrottleOperator is an [Operator] type for [Throttle].
type ThrottleOperator[T, U any] struct {
	ts throttleConfig[T, U]
}

// WithLeading sets Leading option to a given value.
func (op ThrottleOperator[T, U]) WithLeading(v bool) ThrottleOperator[T, U] {
	op.ts.leading = v
	return op
}

// WithTrailing sets Trailing option to a given value.
func (op ThrottleOperator[T, U]) WithTrailing(v bool) ThrottleOperator[T, U] {
	op.ts.trailing = v
	return op
}

// Apply implements the [Operator] interface.
func (op ThrottleOperator[T, U]) Apply(source Observable[T]) Observable[T] {
	return throttleObservable[T, U]{source, op.ts}.Subscribe
}

type throttleObservable[T, U any] struct {
	source Observable[T]
	throttleConfig[T, U]
}

func (ob throttleObservable[T, U]) Subscribe(c Context, o Observer[T]) {
	c, cancel := c.WithCancel()
	o = o.DoOnTermination(cancel)

	var x struct {
		context  atomic.Value
		complete atomic.Bool
		trailing struct {
			sync.Mutex
			value    T
			hasValue atomic.Bool
		}
		worker struct {
			sync.WaitGroup
		}
	}

	x.context.Store(c.Context)

	var doThrottle func(T)

	doThrottle = func(v T) {
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
				if ob.trailing && x.trailing.hasValue.Load() {
					x.trailing.Lock()
					value := x.trailing.value
					x.trailing.hasValue.Store(false)
					x.trailing.Unlock()

					oops := func() {
						if x.context.Swap(sentinel) != sentinel {
							o.Stop(ErrOops)
						}
					}

					Try1(o, Next(value), oops)

					if !x.complete.Load() {
						Try1(doThrottle, value, oops)
					}
				}

				if x.context.CompareAndSwap(w.Context, c.Context) && x.complete.Load() && x.context.CompareAndSwap(c.Context, sentinel) {
					o.Complete()
				}

			case KindComplete:
				if x.context.CompareAndSwap(w.Context, c.Context) && x.complete.Load() && x.context.CompareAndSwap(c.Context, sentinel) {
					o.Complete()
				}

			case KindError, KindStop:
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

	ob.source.Subscribe(c, func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
			x.trailing.Lock()
			x.trailing.value = n.Value
			x.trailing.hasValue.Store(true)
			x.trailing.Unlock()

			if x.context.Load() == c.Context {
				if ob.leading {
					x.trailing.hasValue.Store(false)
					o.Emit(n)
				}

				doThrottle(n.Value)
			}

		case KindComplete:
			x.complete.Store(true)

			if x.context.CompareAndSwap(c.Context, sentinel) {
				o.Emit(n)
			}

		case KindError, KindStop:
			old := x.context.Swap(sentinel)

			cancel()
			x.worker.Wait()

			if old != sentinel {
				o.Emit(n)
			}
		}
	})
}
