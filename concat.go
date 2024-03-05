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

func (obs concatWithObservable[T]) Subscribe(c Context, sink Observer[T]) {
	var observer Observer[T]

	done := c.Done()

	next := resistReentrance(func() {
		if source := obs.Source; source != nil {
			obs.Source = nil
			source.Subscribe(c, observer)
			return
		}

		select {
		default:
		case <-done:
			sink.Error(c.Err())
			return
		}

		if len(obs.Others) == 0 {
			sink.Complete()
			return
		}

		obs1 := obs.Others[0]
		obs.Others = obs.Others[1:]
		obs1.Subscribe(c, observer)
	})

	observer = func(n Notification[T]) {
		if n.Kind != KindComplete {
			sink(n)
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
// by projecting each source value to the same Observable, then flattens it
// into a first-order Observable using ConcatAll.
func ConcatMapTo[T, R any](inner Observable[R]) ConcatMapOperator[T, R] {
	return ConcatMap(func(T) Observable[R] { return inner })
}

// ConcatMap converts the source Observable into a higher-order Observable,
// by projecting each source value to an Observable, then flattens it into
// a first-order Observable using ConcatAll.
func ConcatMap[T, R any](proj func(v T) Observable[R]) ConcatMapOperator[T, R] {
	return ConcatMapOperator[T, R]{
		opts: concatMapConfig[T, R]{
			Project:      proj,
			UseBuffering: false,
		},
	}
}

type concatMapConfig[T, R any] struct {
	Project      func(T) Observable[R]
	UseBuffering bool
}

// ConcatMapOperator is an [Operator] type for [ConcatMap].
type ConcatMapOperator[T, R any] struct {
	opts concatMapConfig[T, R]
}

// WithBuffering turns on source buffering.
// By default, this Operator blocks at every source value until it's flattened.
// With source buffering on, this Operator buffers every source value, which
// might consume a lot of memory over time if the source has lots of values
// emitting faster than concatenating.
func (op ConcatMapOperator[T, R]) WithBuffering() ConcatMapOperator[T, R] {
	op.opts.UseBuffering = true
	return op
}

// Apply implements the Operator interface.
func (op ConcatMapOperator[T, R]) Apply(source Observable[T]) Observable[R] {
	return concatMapObservable[T, R]{source, op.opts}.Subscribe
}

type concatMapObservable[T, R any] struct {
	Source Observable[T]
	concatMapConfig[T, R]
}

func (obs concatMapObservable[T, R]) Subscribe(c Context, sink Observer[R]) {
	if obs.UseBuffering {
		obs.SubscribeWithBuffering(c, sink)
		return
	}

	c, cancel := c.WithCancel()
	sink = sink.OnLastNotification(cancel)

	var noop bool

	obs.Source.Subscribe(c, func(n Notification[T]) {
		if noop {
			return
		}

		switch n.Kind {
		case KindNext:
			if err := obs.Project(n.Value).BlockingSubscribe(c, sink.ElementsOnly); err != nil {
				noop = true
				sink.Error(err)
			}
		case KindError:
			sink.Error(n.Error)
		case KindComplete:
			sink.Complete()
		}
	})
}

func (obs concatMapObservable[T, R]) SubscribeWithBuffering(c Context, sink Observer[R]) {
	c, cancel := c.WithCancel()
	sink = sink.OnLastNotification(cancel)

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

		obs1 := obs.Project(v)
		w, cancelw := c.WithCancel()

		if !x.Context.CompareAndSwap(c.Context, w.Context) { // This fails if x.Context was swapped to sentinel.
			cancelw()
			return
		}

		x.Worker.Add(1)

		obs1.Subscribe(w, func(n Notification[R]) {
			switch n.Kind {
			case KindNext:
				sink(n)

			case KindError, KindComplete:
				defer x.Worker.Done()

				x.Queue.Lock()

				if n.Kind == KindComplete {
					select {
					default:
					case <-w.Done():
						n = Error[R](w.Err())
					}
				}

				cancelw()

				switch n.Kind {
				case KindError:
					old := x.Context.Swap(sentinel)

					x.Queue.Init()
					x.Queue.Unlock()

					if old != sentinel {
						sink(n)
					}

				case KindComplete:
					swapped := x.Context.CompareAndSwap(w.Context, c.Context)
					if swapped && x.Queue.Len() != 0 {
						startWorker()
						return
					}

					x.Queue.Unlock()

					if swapped && x.Complete.Load() && x.Context.CompareAndSwap(c.Context, sentinel) {
						sink.Complete()
					}
				}
			}
		})
	})

	obs.Source.Subscribe(c, func(n Notification[T]) {
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
				sink.Error(n.Error)
			}

		case KindComplete:
			x.Complete.Store(true)

			if x.Context.CompareAndSwap(c.Context, sentinel) {
				sink.Complete()
			}
		}
	})
}
