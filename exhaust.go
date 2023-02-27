package rx

import (
	"context"
	"sync/atomic"
)

// ExhaustAll flattens a higher-order Observable into a first-order Observable
// by dropping inner Observables while the previous inner Observable has not
// yet completed.
func ExhaustAll[_ Observable[T], T any]() Operator[Observable[T], T] {
	return exhaustMap(identity[Observable[T]])
}

// ExhaustMap converts the source Observable into a higher-order Observable,
// by projecting each source value to an Observable, then flattens it into
// a first-order Observable using ExhaustAll.
func ExhaustMap[T, R any](proj func(v T) Observable[R]) Operator[T, R] {
	if proj == nil {
		panic("proj == nil")
	}

	return exhaustMap(proj)
}

// ExhaustMapTo converts the source Observable into a higher-order Observable,
// by projecting each source value to the same Observable, then flattens it
// into a first-order Observable using ExhaustAll.
func ExhaustMapTo[T, R any](inner Observable[R]) Operator[T, R] {
	return exhaustMap(func(T) Observable[R] { return inner })
}

func exhaustMap[T, R any](proj func(v T) Observable[R]) Operator[T, R] {
	return NewOperator(
		func(source Observable[T]) Observable[R] {
			return exhaustMapObservable[T, R]{source, proj}.Subscribe
		},
	)
}

type exhaustMapObservable[T, R any] struct {
	Source  Observable[T]
	Project func(T) Observable[R]
}

func (obs exhaustMapObservable[T, R]) Subscribe(ctx context.Context, sink Observer[R]) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel).WithMutex()

	var workers atomic.Uint32

	workers.Store(1)

	observer := func(n Notification[R]) {
		if n.HasValue || n.HasError || workers.Add(^uint32(0)) == 0 {
			sink(n)
		}
	}

	obs.Source.Subscribe(ctx, func(n Notification[T]) {
		switch {
		case n.HasValue:
			if workers.CompareAndSwap(1, 2) {
				obs.Project(n.Value).Subscribe(ctx, observer)
			}
		case n.HasError:
			sink.Error(n.Error)
		default:
			if workers.Add(^uint32(0)) == 0 {
				sink.Complete()
			}
		}
	})
}
