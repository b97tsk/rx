package rx

import (
	"sync"
	"sync/atomic"

	"github.com/b97tsk/rx/internal/queue"
)

// Merge creates an Observable that concurrently emits all values from every
// given input Observable.
func Merge[T any](some ...Observable[T]) Observable[T] {
	if len(some) == 0 {
		return Empty[T]()
	}

	return observables[T](some).Merge
}

// MergeWith applies [Merge] to the source Observable along with some other
// Observables to create a first-order Observable.
func MergeWith[T any](some ...Observable[T]) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return observables[T](append([]Observable[T]{source}, some...)).Merge
		},
	)
}

func (some observables[T]) Merge(c Context, sink Observer[T]) {
	c, cancel := c.WithCancel()
	sink = sink.OnLastNotification(cancel).Serialized()

	var workers atomic.Uint32

	workers.Store(uint32(len(some)))

	observer := func(n Notification[T]) {
		if n.Kind != KindComplete || workers.Add(^uint32(0)) == 0 {
			sink(n)
		}
	}

	for _, obs := range some {
		c.Go(func() { obs.Subscribe(c, observer) })
	}
}

// MergeAll flattens a higher-order Observable into a first-order Observable
// which concurrently delivers all values that are emitted on the inner
// Observables.
func MergeAll[_ Observable[T], T any]() MergeMapOperator[Observable[T], T] {
	return MergeMap(identity[Observable[T]])
}

// MergeMapTo converts the source Observable into a higher-order Observable,
// by projecting each source value to the same Observable, then flattens it
// into a first-order Observable using MergeAll.
func MergeMapTo[T, R any](inner Observable[R]) MergeMapOperator[T, R] {
	return MergeMap(func(T) Observable[R] { return inner })
}

// MergeMap converts the source Observable into a higher-order Observable,
// by projecting each source value to an Observable, then flattens it into
// a first-order Observable using MergeAll.
func MergeMap[T, R any](proj func(v T) Observable[R]) MergeMapOperator[T, R] {
	return MergeMapOperator[T, R]{
		opts: mergeMapConfig[T, R]{
			Project:      proj,
			Concurrency:  -1,
			UseBuffering: false,
		},
	}
}

type mergeMapConfig[T, R any] struct {
	Project      func(T) Observable[R]
	Concurrency  int
	UseBuffering bool
}

// MergeMapOperator is an [Operator] type for [MergeMap].
type MergeMapOperator[T, R any] struct {
	opts mergeMapConfig[T, R]
}

// WithBuffering turns on source buffering.
// By default, this Operator might block the source due to concurrency limit.
// With source buffering on, this Operator buffers every source value, which
// might consume a lot of memory over time if the source has lots of values
// emitting faster than merging.
func (op MergeMapOperator[T, R]) WithBuffering() MergeMapOperator[T, R] {
	op.opts.UseBuffering = true
	return op
}

// WithConcurrency sets Concurrency option to a given value.
// It must not be zero. The default value is -1 (unlimited).
func (op MergeMapOperator[T, R]) WithConcurrency(n int) MergeMapOperator[T, R] {
	op.opts.Concurrency = n
	return op
}

// Apply implements the Operator interface.
func (op MergeMapOperator[T, R]) Apply(source Observable[T]) Observable[R] {
	if op.opts.Concurrency == 0 {
		return Oops[R]("MergeMap: Concurrency == 0")
	}

	return mergeMapObservable[T, R]{source, op.opts}.Subscribe
}

type mergeMapObservable[T, R any] struct {
	Source Observable[T]
	mergeMapConfig[T, R]
}

func (obs mergeMapObservable[T, R]) Subscribe(c Context, sink Observer[R]) {
	if obs.UseBuffering {
		obs.SubscribeWithBuffering(c, sink)
		return
	}

	c, cancel := c.WithCancel()
	sink = sink.OnLastNotification(cancel).Serialized()

	var x struct {
		sync.Mutex
		sync.Cond
		Workers  int
		Complete bool
		HasError bool
	}

	x.Cond.L = &x.Mutex

	var noop bool

	obs.Source.Subscribe(c, func(n Notification[T]) {
		if noop {
			return
		}

		switch n.Kind {
		case KindNext:
			x.Lock()

			for x.Workers == obs.Concurrency && !x.HasError {
				x.Wait()
			}

			if x.HasError {
				noop = true
				x.Unlock()
				return
			}

			obs1 := Try11(obs.Project, n.Value, func() {
				defer x.Unlock()
				noop = true
				x.HasError = true
				sink.Error(ErrOops)
			})

			x.Workers++
			x.Unlock()

			c.Go(func() {
				obs1.Subscribe(c, func(n Notification[R]) {
					switch n.Kind {
					case KindNext:
						sink(n)

					case KindError:
						x.Lock()
						x.Workers--
						x.HasError = true
						x.Unlock()
						x.Signal()

						sink(n)

					case KindComplete:
						x.Lock()

						x.Workers--

						if x.Workers == 0 && x.Complete && !x.HasError {
							x.Unlock()
							sink(n)
							return
						}

						x.Unlock()
						x.Signal()
					}
				})
			})

		case KindError:
			sink.Error(n.Error)

		case KindComplete:
			x.Lock()

			x.Complete = true

			if x.Workers == 0 && !x.HasError {
				x.Unlock()
				sink.Complete()
				return
			}

			x.Unlock()
		}
	})
}

func (obs mergeMapObservable[T, R]) SubscribeWithBuffering(c Context, sink Observer[R]) {
	c, cancel := c.WithCancel()
	sink = sink.OnLastNotification(cancel).Serialized()

	var x struct {
		sync.Mutex
		Queue    queue.Queue[T]
		Workers  int
		Complete bool
		HasError bool
	}

	var startWorker func()

	startWorker = func() {
		obs1 := Try11(obs.Project, x.Queue.Pop(), func() {
			defer x.Unlock()
			x.Queue.Init()
			x.HasError = true
			sink.Error(ErrOops)
		})

		x.Workers++
		x.Unlock()

		c.Go(func() {
			obs1.Subscribe(c, func(n Notification[R]) {
				switch n.Kind {
				case KindNext:
					sink(n)

				case KindError:
					x.Lock()
					x.Queue.Init()
					x.Workers--
					x.HasError = true
					x.Unlock()

					sink(n)

				case KindComplete:
					x.Lock()

					x.Workers--

					if x.Queue.Len() != 0 {
						startWorker()
						return
					}

					if x.Workers == 0 && x.Complete && !x.HasError {
						x.Unlock()
						sink(n)
						return
					}

					x.Unlock()
				}
			})
		})
	}

	var noop bool

	obs.Source.Subscribe(c, func(n Notification[T]) {
		if noop {
			return
		}

		switch n.Kind {
		case KindNext:
			x.Lock()

			if x.HasError {
				noop = true
				x.Unlock()
				return
			}

			x.Queue.Push(n.Value)

			if x.Workers != obs.Concurrency {
				startWorker()
				return
			}

			x.Unlock()

		case KindError:
			sink.Error(n.Error)

		case KindComplete:
			x.Lock()

			x.Complete = true

			if x.Workers == 0 && !x.HasError {
				x.Unlock()
				sink.Complete()
				return
			}

			x.Unlock()
		}
	})
}
