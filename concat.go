package rx

import (
	"sync"
	"sync/atomic"

	"github.com/b97tsk/rx/internal/queue"
)

// Concat creates an Observable that concatenates multiple Observables
// together by sequentially emitting their values, one Observable after
// the other.
func Concat[T any](some ...Observable[T]) Observable[T] {
	if len(some) == 0 {
		return Empty[T]()
	}

	return concatWithObservable[T]{others: some}.Subscribe
}

// ConcatWith concatenates the source Observable and some other Observables
// together to create an Observable that sequentially emits their values,
// one Observable after the other.
func ConcatWith[T any](some ...Observable[T]) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return concatWithObservable[T]{source, some}.Subscribe
		},
	)
}

type concatWithObservable[T any] struct {
	source Observable[T]
	others []Observable[T]
}

func (ob concatWithObservable[T]) Subscribe(c Context, o Observer[T]) {
	var observer Observer[T]

	done := c.Done()

	next := resistReentrance(func() {
		if source := ob.source; source != nil {
			ob.source = nil
			source.Subscribe(c, observer)
			return
		}

		select {
		default:
		case <-done:
			o.Error(c.Cause())
			return
		}

		if len(ob.others) == 0 {
			o.Complete()
			return
		}

		obs := ob.others[0]
		ob.others = ob.others[1:]
		obs.Subscribe(c, observer)
	})

	observer = func(n Notification[T]) {
		if n.Kind != KindComplete {
			o.Emit(n)
			return
		}

		next()
	}

	next()
}

// ConcatAll flattens a higher-order Observable into a first-order Observable
// by concatenating the inner Observables in order.
func ConcatAll[_ Observable[T], T any]() ConcatMapOperator[Observable[T], T] {
	return ConcatMap(identity[Observable[T]])
}

// ConcatMapTo converts the source Observable into a higher-order Observable,
// by mapping each source value to the same Observable, then flattens it into
// a first-order Observable using ConcatAll.
func ConcatMapTo[T, R any](inner Observable[R]) ConcatMapOperator[T, R] {
	return ConcatMap(func(T) Observable[R] { return inner })
}

// ConcatMap converts the source Observable into a higher-order Observable,
// by mapping each source value to an Observable, then flattens it into
// a first-order Observable using ConcatAll.
func ConcatMap[T, R any](mapping func(v T) Observable[R]) ConcatMapOperator[T, R] {
	return ConcatMapOperator[T, R]{
		ts: concatMapConfig[T, R]{
			mapping:      mapping,
			useBuffering: false,
		},
	}
}

type concatMapConfig[T, R any] struct {
	mapping      func(T) Observable[R]
	useBuffering bool
}

// ConcatMapOperator is an [Operator] type for [ConcatMap].
type ConcatMapOperator[T, R any] struct {
	ts concatMapConfig[T, R]
}

// WithBuffering turns on source buffering.
// By default, this Operator blocks at every source value until it's flattened.
// With source buffering on, this Operator buffers every source value, which
// might consume a lot of memory over time if the source has lots of values
// emitting faster than concatenating.
func (op ConcatMapOperator[T, R]) WithBuffering() ConcatMapOperator[T, R] {
	op.ts.useBuffering = true
	return op
}

// Apply implements the Operator interface.
func (op ConcatMapOperator[T, R]) Apply(source Observable[T]) Observable[R] {
	return concatMapObservable[T, R]{source, op.ts}.Subscribe
}

type concatMapObservable[T, R any] struct {
	source Observable[T]
	concatMapConfig[T, R]
}

func (ob concatMapObservable[T, R]) Subscribe(c Context, o Observer[R]) {
	if ob.useBuffering {
		ob.SubscribeWithBuffering(c, o)
		return
	}

	c, cancel := c.WithCancel()
	o = o.DoOnTermination(cancel)

	var noop bool

	ob.source.Subscribe(c, func(n Notification[T]) {
		if noop {
			return
		}

		switch n.Kind {
		case KindNext:
			if err := ob.mapping(n.Value).BlockingSubscribe(c, o.ElementsOnly); err != nil {
				noop = true
				o.Error(err)
			}
		case KindError:
			o.Error(n.Error)
		case KindComplete:
			o.Complete()
		}
	})
}

func (ob concatMapObservable[T, R]) SubscribeWithBuffering(c Context, o Observer[R]) {
	c, cancel := c.WithCancel()
	o = o.DoOnTermination(cancel)

	var x struct {
		context  atomic.Value
		complete atomic.Bool
		buffer   struct {
			sync.Mutex
			queue.Queue[T]
		}
		worker struct {
			sync.WaitGroup
		}
	}

	x.context.Store(c.Context)

	var startWorker func()

	startWorker = resistReentrance(func() {
		v := x.buffer.Pop()

		x.buffer.Unlock()

		obs := ob.mapping(v)
		w, cancelw := c.WithCancel()

		if !x.context.CompareAndSwap(c.Context, w.Context) { // This fails if x.Context was swapped to sentinel.
			cancelw()
			return
		}

		x.worker.Add(1)

		obs.Subscribe(w, func(n Notification[R]) {
			switch n.Kind {
			case KindNext:
				o.Emit(n)

			case KindError, KindComplete:
				defer x.worker.Done()

				x.buffer.Lock()

				if n.Kind == KindComplete {
					select {
					default:
					case <-w.Done():
						n = Error[R](w.Cause())
					}
				}

				cancelw()

				switch n.Kind {
				case KindError:
					old := x.context.Swap(sentinel)

					x.buffer.Init()
					x.buffer.Unlock()

					if old != sentinel {
						o.Emit(n)
					}

				case KindComplete:
					swapped := x.context.CompareAndSwap(w.Context, c.Context)
					if swapped && x.buffer.Len() != 0 {
						startWorker()
						return
					}

					x.buffer.Unlock()

					if swapped && x.complete.Load() && x.context.CompareAndSwap(c.Context, sentinel) {
						o.Complete()
					}
				}
			}
		})
	})

	ob.source.Subscribe(c, func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
			x.buffer.Lock()

			ctx := x.context.Load()
			if ctx != sentinel {
				x.buffer.Push(n.Value)
			}

			if ctx == c.Context {
				startWorker()
				return
			}

			x.buffer.Unlock()

		case KindError:
			old := x.context.Swap(sentinel)

			cancel()
			x.worker.Wait()
			x.buffer.Init()

			if old != sentinel {
				o.Error(n.Error)
			}

		case KindComplete:
			x.complete.Store(true)

			if x.context.CompareAndSwap(c.Context, sentinel) {
				o.Complete()
			}
		}
	})
}
