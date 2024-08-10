package rx

import (
	"sync"
	"sync/atomic"

	"github.com/b97tsk/rx/internal/queue"
)

// Merge creates an [Observable] that concurrently emits all values from
// every given input [Observable].
func Merge[T any](some ...Observable[T]) Observable[T] {
	if len(some) == 0 {
		return Empty[T]()
	}

	return mergeWithObservable[T]{others: some}.Subscribe
}

// MergeWith merges the source [Observable] and some other Observables
// together to create an [Observable] that concurrently emits all values
// from the source and every given input [Observable].
func MergeWith[T any](some ...Observable[T]) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return mergeWithObservable[T]{source, some}.Subscribe
		},
	)
}

type mergeWithObservable[T any] struct {
	source Observable[T]
	others []Observable[T]
}

func (ob mergeWithObservable[T]) Subscribe(c Context, o Observer[T]) {
	c, o = Synchronize(c, o)

	var num atomic.Uint32

	num.Store(uint32(ob.numObservables()))

	worker := func(n Notification[T]) {
		if n.Kind != KindComplete || num.Add(^uint32(0)) == 0 {
			o.Emit(n)
		}
	}

	done := c.Done()

	if ob.source != nil {
		ob.source.Subscribe(c, worker)

		select {
		default:
		case <-done:
			o.Stop(c.Cause())
			return
		}
	}

	for _, obs := range ob.others {
		obs.Subscribe(c, worker)

		select {
		default:
		case <-done:
			o.Stop(c.Cause())
			return
		}
	}
}

func (ob mergeWithObservable[T]) numObservables() int {
	n := len(ob.others)

	if ob.source != nil {
		n++
	}

	return n
}

// MergeAll flattens a higher-order [Observable] into a first-order
// [Observable] which concurrently delivers all values that are emitted
// on the inner Observables.
func MergeAll[_ Observable[T], T any]() MergeMapOperator[Observable[T], T] {
	return MergeMap(identity[Observable[T]])
}

// MergeMapTo converts the source [Observable] into a higher-order
// [Observable], by mapping each source value to the same [Observable],
// then flattens it into a first-order [Observable] using [MergeAll].
func MergeMapTo[T, R any](inner Observable[R]) MergeMapOperator[T, R] {
	return MergeMap(func(T) Observable[R] { return inner })
}

// MergeMap converts the source [Observable] into a higher-order [Observable],
// by mapping each source value to an [Observable], then flattens it into
// a first-order [Observable] using [MergeAll].
func MergeMap[T, R any](mapping func(v T) Observable[R]) MergeMapOperator[T, R] {
	return MergeMapOperator[T, R]{
		ts: mergeMapConfig[T, R]{
			mapping:      mapping,
			concurrency:  -1,
			useBuffering: false,
		},
	}
}

type mergeMapConfig[T, R any] struct {
	mapping      func(T) Observable[R]
	concurrency  int
	useBuffering bool
}

// MergeMapOperator is an [Operator] type for [MergeMap].
type MergeMapOperator[T, R any] struct {
	ts mergeMapConfig[T, R]
}

// WithBuffering turns on source buffering.
// By default, this Operator might block the source due to concurrency limit.
// With source buffering on, this Operator buffers every source value, which
// might consume a lot of memory over time if the source has lots of values
// emitting faster than merging.
func (op MergeMapOperator[T, R]) WithBuffering() MergeMapOperator[T, R] {
	op.ts.useBuffering = true
	return op
}

// WithConcurrency sets Concurrency option to a given value.
// It must not be zero. The default value is -1 (unlimited).
func (op MergeMapOperator[T, R]) WithConcurrency(n int) MergeMapOperator[T, R] {
	op.ts.concurrency = n
	return op
}

// Apply implements the [Operator] interface.
func (op MergeMapOperator[T, R]) Apply(source Observable[T]) Observable[R] {
	if op.ts.concurrency == 0 {
		return Oops[R]("MergeMap: Concurrency == 0")
	}

	return mergeMapObservable[T, R]{source, op.ts}.Subscribe
}

type mergeMapObservable[T, R any] struct {
	source Observable[T]
	mergeMapConfig[T, R]
}

func (ob mergeMapObservable[T, R]) Subscribe(c Context, o Observer[R]) {
	if ob.useBuffering {
		ob.SubscribeWithBuffering(c, o)
		return
	}

	c, o = Synchronize(c, o)

	var x struct {
		mu sync.Mutex
		co sync.Cond

		workers  int
		complete bool
		hasError bool
	}

	x.co.L = &x.mu

	worker := func(n Notification[R]) {
		switch n.Kind {
		case KindNext:
			o.Emit(n)

		case KindComplete:
			x.mu.Lock()

			x.workers--

			if x.workers == 0 && x.complete && !x.hasError {
				x.mu.Unlock()
				o.Emit(n)
				return
			}

			x.mu.Unlock()
			x.co.Signal()

		case KindError, KindStop:
			x.mu.Lock()
			x.workers--
			x.hasError = true
			x.mu.Unlock()
			x.co.Signal()
			o.Emit(n)
		}
	}

	var noop bool

	ob.source.Subscribe(c, func(n Notification[T]) {
		if noop {
			return
		}

		switch n.Kind {
		case KindNext:
			x.mu.Lock()

			for x.workers == ob.concurrency && !x.hasError {
				x.co.Wait()
			}

			if x.hasError {
				noop = true
				x.mu.Unlock()
				return
			}

			obs := Try11(ob.mapping, n.Value, func() {
				noop = true
				x.hasError = true
				x.mu.Unlock()
				o.Stop(ErrOops)
			})

			x.workers++
			x.mu.Unlock()

			obs.Subscribe(c, worker)

		case KindComplete:
			x.mu.Lock()

			x.complete = true

			if x.workers == 0 && !x.hasError {
				x.mu.Unlock()
				o.Complete()
				return
			}

			x.mu.Unlock()

		case KindError:
			o.Error(n.Error)

		case KindStop:
			o.Stop(n.Error)
		}
	})
}

func (ob mergeMapObservable[T, R]) SubscribeWithBuffering(c Context, o Observer[R]) {
	c, o = Synchronize(c, o)

	var x struct {
		mu sync.Mutex

		buffer   queue.Queue[T]
		workers  int
		complete bool
		hasError bool

		startWorkerPool sync.Pool
	}

	startWorker := func() {
		if v := x.startWorkerPool.Get(); v != nil {
			v.(func())()
			return
		}

		var startWorker func()

		worker := func(n Notification[R]) {
			switch n.Kind {
			case KindNext:
				o.Emit(n)

			case KindComplete:
				x.mu.Lock()

				x.workers--

				if x.buffer.Len() != 0 {
					startWorker()
					return
				}

				if x.workers == 0 && x.complete && !x.hasError {
					x.mu.Unlock()
					o.Emit(n)
					return
				}

				x.startWorkerPool.Put(startWorker)

				x.mu.Unlock()

			case KindError, KindStop:
				x.mu.Lock()
				x.buffer.Init()
				x.workers--
				x.hasError = true
				x.mu.Unlock()
				o.Emit(n)
			}
		}

		startWorker = resistReentrance(func() {
			obs := Try11(ob.mapping, x.buffer.Pop(), func() {
				x.buffer.Init()
				x.hasError = true
				x.mu.Unlock()
				o.Stop(ErrOops)
			})

			x.workers++
			x.mu.Unlock()

			obs.Subscribe(c, worker)
		})

		startWorker()
	}

	var noop bool

	ob.source.Subscribe(c, func(n Notification[T]) {
		if noop {
			return
		}

		switch n.Kind {
		case KindNext:
			x.mu.Lock()

			if x.hasError {
				noop = true
				x.mu.Unlock()
				return
			}

			x.buffer.Push(n.Value)

			if x.workers != ob.concurrency {
				startWorker()
				return
			}

			x.mu.Unlock()

		case KindComplete:
			x.mu.Lock()

			x.complete = true

			if x.workers == 0 && !x.hasError {
				x.mu.Unlock()
				o.Complete()
				return
			}

			x.mu.Unlock()

		case KindError:
			o.Error(n.Error)

		case KindStop:
			o.Stop(n.Error)
		}
	})
}
