package rx

import (
	"context"

	"github.com/b97tsk/rx/internal/atomic"
)

// Merge creates an Observable that concurrently emits all values from every
// given input Observable.
func Merge[T any](some ...Observable[T]) Observable[T] {
	if len(some) == 0 {
		return Empty[T]()
	}

	return observables[T](some).Merge
}

// MergeWith applies Merge to the source Observable along with some other
// Observables to create a first-order Observable, then mirrors the resulting
// Observable.
func MergeWith[T any](some ...Observable[T]) Operator[T, T] {
	return AsOperator(
		func(source Observable[T]) Observable[T] {
			return observables[T](append([]Observable[T]{source}, some...)).Merge
		},
	)
}

func (some observables[T]) Merge(ctx context.Context, sink Observer[T]) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel).Mutex()

	workers := atomic.FromUint32(uint32(len(some)))

	observer := func(n Notification[T]) {
		if n.HasValue || n.HasError || workers.Sub(1) == 0 {
			sink(n)
		}
	}

	for _, obs := range some {
		go obs.Subscribe(ctx, observer)
	}
}
