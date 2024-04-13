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

	return concatWithObservable[T]{Others: some}.Subscribe
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
	Source Observable[T]
	Others []Observable[T]
}

func (ob concatWithObservable[T]) Subscribe(c Context, o Observer[T]) {
	var observer Observer[T]

	done := c.Done()

	next := resistReentrance(func() {
		if source := ob.Source; source != nil {
			ob.Source = nil
			source.Subscribe(c, observer)
			return
		}

		select {
		default:
		case <-done:
			o.Error(c.Cause())
			return
		}

		if len(ob.Others) == 0 {
			o.Complete()
			return
		}

		obs := ob.Others[0]
		ob.Others = ob.Others[1:]
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
			Mapping:      mapping,
			UseBuffering: false,
		},
	}
}

type concatMapConfig[T, R any] struct {
	Mapping      func(T) Observable[R]
	UseBuffering bool
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
	op.ts.UseBuffering = true
	return op
}

// Apply implements the Operator interface.
func (op ConcatMapOperator[T, R]) Apply(source Observable[T]) Observable[R] {
	return concatMapObservable[T, R]{source, op.ts}.Subscribe
}

type concatMapObservable[T, R any] struct {
	Source Observable[T]
	concatMapConfig[T, R]
}

func (ob concatMapObservable[T, R]) Subscribe(c Context, o Observer[R]) {
	if ob.UseBuffering {
		ob.SubscribeWithBuffering(c, o)
		return
	}

	c, cancel := c.WithCancel()
	o = o.DoOnTermination(cancel)

	var noop bool

	ob.Source.Subscribe(c, func(n Notification[T]) {
		if noop {
			return
		}

		switch n.Kind {
		case KindNext:
			if err := ob.Mapping(n.Value).BlockingSubscribe(c, o.ElementsOnly); err != nil {
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
		Context  atomic.Value
		Complete atomic.Bool
		Queue    struct {
			sync.Mutex
			queue.Queue[T]
		}
		Worker struct {
			sync.WaitGroup
		}
	}

	x.Context.Store(c.Context)

	var startWorker func()

	startWorker = resistReentrance(func() {
		v := x.Queue.Pop()

		x.Queue.Unlock()

		obs := ob.Mapping(v)
		w, cancelw := c.WithCancel()

		if !x.Context.CompareAndSwap(c.Context, w.Context) { // This fails if x.Context was swapped to sentinel.
			cancelw()
			return
		}

		x.Worker.Add(1)

		obs.Subscribe(w, func(n Notification[R]) {
			switch n.Kind {
			case KindNext:
				o.Emit(n)

			case KindError, KindComplete:
				defer x.Worker.Done()

				x.Queue.Lock()

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
					old := x.Context.Swap(sentinel)

					x.Queue.Init()
					x.Queue.Unlock()

					if old != sentinel {
						o.Emit(n)
					}

				case KindComplete:
					swapped := x.Context.CompareAndSwap(w.Context, c.Context)
					if swapped && x.Queue.Len() != 0 {
						startWorker()
						return
					}

					x.Queue.Unlock()

					if swapped && x.Complete.Load() && x.Context.CompareAndSwap(c.Context, sentinel) {
						o.Complete()
					}
				}
			}
		})
	})

	ob.Source.Subscribe(c, func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
			x.Queue.Lock()

			ctx := x.Context.Load()
			if ctx != sentinel {
				x.Queue.Push(n.Value)
			}

			if ctx == c.Context {
				startWorker()
				return
			}

			x.Queue.Unlock()

		case KindError:
			old := x.Context.Swap(sentinel)

			cancel()
			x.Worker.Wait()
			x.Queue.Init()

			if old != sentinel {
				o.Error(n.Error)
			}

		case KindComplete:
			x.Complete.Store(true)

			if x.Context.CompareAndSwap(c.Context, sentinel) {
				o.Complete()
			}
		}
	})
}
