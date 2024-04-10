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

	return mergeWithObservable[T]{Others: some}.Subscribe
}

// MergeWith merges the source Observable and some other Observables together
// to create an Observable that concurrently emits all values from the source
// and every given input Observable.
func MergeWith[T any](some ...Observable[T]) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return mergeWithObservable[T]{source, some}.Subscribe
		},
	)
}

type mergeWithObservable[T any] struct {
	Source Observable[T]
	Others []Observable[T]
}

func (obs mergeWithObservable[T]) Subscribe(c Context, o Observer[T]) {
	c, o = Serialize(c, o)

	var num atomic.Uint32

	num.Store(uint32(obs.numObservables()))

	worker := func(n Notification[T]) {
		if n.Kind != KindComplete || num.Add(^uint32(0)) == 0 {
			o.Emit(n)
		}
	}

	done := c.Done()

	if obs.Source != nil {
		obs.Source.Subscribe(c, worker)

		select {
		default:
		case <-done:
			o.Error(c.Err())
			return
		}
	}

	for _, obs1 := range obs.Others {
		obs1.Subscribe(c, worker)

		select {
		default:
		case <-done:
			o.Error(c.Err())
			return
		}
	}
}

func (obs mergeWithObservable[T]) numObservables() int {
	n := len(obs.Others)

	if obs.Source != nil {
		n++
	}

	return n
}

// MergeAll flattens a higher-order Observable into a first-order Observable
// which concurrently delivers all values that are emitted on the inner
// Observables.
func MergeAll[_ Observable[T], T any]() MergeMapOperator[Observable[T], T] {
	return MergeMap(identity[Observable[T]])
}

// MergeMapTo converts the source Observable into a higher-order Observable,
// by mapping each source value to the same Observable, then flattens it into
// a first-order Observable using MergeAll.
func MergeMapTo[T, R any](inner Observable[R]) MergeMapOperator[T, R] {
	return MergeMap(func(T) Observable[R] { return inner })
}

// MergeMap converts the source Observable into a higher-order Observable,
// by mapping each source value to an Observable, then flattens it into
// a first-order Observable using MergeAll.
func MergeMap[T, R any](mapping func(v T) Observable[R]) MergeMapOperator[T, R] {
	return MergeMapOperator[T, R]{
		ts: mergeMapConfig[T, R]{
			Mapping:      mapping,
			Concurrency:  -1,
			UseBuffering: false,
		},
	}
}

type mergeMapConfig[T, R any] struct {
	Mapping      func(T) Observable[R]
	Concurrency  int
	UseBuffering bool
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
	op.ts.UseBuffering = true
	return op
}

// WithConcurrency sets Concurrency option to a given value.
// It must not be zero. The default value is -1 (unlimited).
func (op MergeMapOperator[T, R]) WithConcurrency(n int) MergeMapOperator[T, R] {
	op.ts.Concurrency = n
	return op
}

// Apply implements the Operator interface.
func (op MergeMapOperator[T, R]) Apply(source Observable[T]) Observable[R] {
	if op.ts.Concurrency == 0 {
		return Oops[R]("MergeMap: Concurrency == 0")
	}

	return mergeMapObservable[T, R]{source, op.ts}.Subscribe
}

type mergeMapObservable[T, R any] struct {
	Source Observable[T]
	mergeMapConfig[T, R]
}

func (obs mergeMapObservable[T, R]) Subscribe(c Context, o Observer[R]) {
	if obs.UseBuffering {
		obs.SubscribeWithBuffering(c, o)
		return
	}

	c, o = Serialize(c, o)

	var x struct {
		sync.Mutex
		sync.Cond
		Workers  int
		Complete bool
		HasError bool
	}

	x.Cond.L = &x.Mutex

	worker := func(n Notification[R]) {
		switch n.Kind {
		case KindNext:
			o.Emit(n)

		case KindError:
			x.Lock()
			x.Workers--
			x.HasError = true
			x.Unlock()
			x.Signal()

			o.Emit(n)

		case KindComplete:
			x.Lock()

			x.Workers--

			if x.Workers == 0 && x.Complete && !x.HasError {
				x.Unlock()
				o.Emit(n)
				return
			}

			x.Unlock()
			x.Signal()
		}
	}

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

			obs1 := Try11(obs.Mapping, n.Value, func() {
				defer x.Unlock()
				noop = true
				x.HasError = true
				o.Error(ErrOops)
			})

			x.Workers++
			x.Unlock()

			obs1.Subscribe(c, worker)

		case KindError:
			o.Error(n.Error)

		case KindComplete:
			x.Lock()

			x.Complete = true

			if x.Workers == 0 && !x.HasError {
				x.Unlock()
				o.Complete()
				return
			}

			x.Unlock()
		}
	})
}

func (obs mergeMapObservable[T, R]) SubscribeWithBuffering(c Context, o Observer[R]) {
	c, o = Serialize(c, o)

	var x struct {
		sync.Mutex

		Queue    queue.Queue[T]
		Workers  int
		Complete bool
		HasError bool

		StartWorkerPool sync.Pool
	}

	startWorker := func() {
		if v := x.StartWorkerPool.Get(); v != nil {
			(*v.(*func()))()
			return
		}

		var startWorker func()

		worker := func(n Notification[R]) {
			switch n.Kind {
			case KindNext:
				o.Emit(n)

			case KindError:
				x.Lock()
				x.Queue.Init()
				x.Workers--
				x.HasError = true
				x.Unlock()

				o.Emit(n)

			case KindComplete:
				x.Lock()

				x.Workers--

				if x.Queue.Len() != 0 {
					startWorker()
					return
				}

				if x.Workers == 0 && x.Complete && !x.HasError {
					x.Unlock()
					o.Emit(n)
					return
				}

				x.StartWorkerPool.Put(&startWorker)

				x.Unlock()
			}
		}

		startWorker = resistReentrance(func() {
			obs1 := Try11(obs.Mapping, x.Queue.Pop(), func() {
				defer x.Unlock()
				x.Queue.Init()
				x.HasError = true
				o.Error(ErrOops)
			})

			x.Workers++
			x.Unlock()

			obs1.Subscribe(c, worker)
		})

		startWorker()
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
			o.Error(n.Error)

		case KindComplete:
			x.Lock()

			x.Complete = true

			if x.Workers == 0 && !x.HasError {
				x.Unlock()
				o.Complete()
				return
			}

			x.Unlock()
		}
	})
}
