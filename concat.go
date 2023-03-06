package rx

import (
	"context"
	"sync"

	"github.com/b97tsk/rx/internal/queue"
)

// Concat creates an Observable that concatenates multiple Observables
// together by sequentially emitting their values, one Observable after
// the other.
func Concat[T any](some ...Observable[T]) Observable[T] {
	if len(some) == 0 {
		return Empty[T]()
	}

	return observables[T](some).Concat
}

// ConcatWith applies [Concat] to the source Observable along with some other
// Observables to create a first-order Observable.
func ConcatWith[T any](some ...Observable[T]) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return observables[T](append([]Observable[T]{source}, some...)).Concat
		},
	)
}

func (some observables[T]) Concat(ctx context.Context, sink Observer[T]) {
	var observer Observer[T]

	subscribeToNext := resistReentry(func() {
		if len(some) == 0 {
			sink.Complete()
			return
		}

		if err := getErr(ctx); err != nil {
			sink.Error(err)
			return
		}

		obs := some[0]
		some = some[1:]

		obs.Subscribe(ctx, observer)
	})

	observer = func(n Notification[T]) {
		if n.HasValue || n.HasError {
			sink(n)
			return
		}

		subscribeToNext()
	}

	subscribeToNext()
}

// ConcatAll flattens a higher-order Observable into a first-order Observable
// by concatenating the inner Observables in order.
func ConcatAll[_ Observable[T], T any]() Operator[Observable[T], T] {
	return concatMap(identity[Observable[T]])
}

// ConcatMap converts the source Observable into a higher-order Observable,
// by projecting each source value to an Observable, then flattens it into
// a first-order Observable using ConcatAll.
func ConcatMap[T, R any](proj func(v T) Observable[R]) Operator[T, R] {
	if proj == nil {
		panic("proj == nil")
	}

	return concatMap(proj)
}

// ConcatMapTo converts the source Observable into a higher-order Observable,
// by projecting each source value to the same Observable, then flattens it
// into a first-order Observable using ConcatAll.
func ConcatMapTo[T, R any](inner Observable[R]) Operator[T, R] {
	return concatMap(func(T) Observable[R] { return inner })
}

func concatMap[T, R any](proj func(v T) Observable[R]) Operator[T, R] {
	return NewOperator(
		func(source Observable[T]) Observable[R] {
			return concatMapObservable[T, R]{source, proj}.Subscribe
		},
	)
}

type concatMapObservable[T, R any] struct {
	Source  Observable[T]
	Project func(T) Observable[R]
}

func (obs concatMapObservable[T, R]) Subscribe(ctx context.Context, sink Observer[R]) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.OnLastNotification(cancel).WithMutex()

	var x struct {
		sync.Mutex
		Queue    queue.Queue[T]
		Working  bool
		Complete bool
	}

	var observer Observer[R]

	subscribeToNext := resistReentry(func() {
		x.Lock()

		if x.Queue.Len() == 0 {
			x.Working = false

			if x.Complete {
				sink.Complete()
			}

			x.Unlock()

			return
		}

		v := x.Queue.Pop()

		x.Unlock()

		obs.Project(v).Subscribe(ctx, observer)
	})

	observer = func(n Notification[R]) {
		if n.HasValue || n.HasError {
			sink(n)
			return
		}

		if err := getErr(ctx); err != nil {
			sink.Error(err)
			return
		}

		subscribeToNext()
	}

	obs.Source.Subscribe(ctx, func(n Notification[T]) {
		switch {
		case n.HasValue:
			x.Lock()

			x.Queue.Push(n.Value)

			var subscribe bool

			if !x.Working {
				x.Working = true
				subscribe = true
			}

			x.Unlock()

			if subscribe {
				subscribeToNext()
			}

		case n.HasError:
			sink.Error(n.Error)

		default:
			x.Lock()

			x.Complete = true

			if !x.Working {
				sink.Complete()
			}

			x.Unlock()
		}
	})
}
