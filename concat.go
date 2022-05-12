package rx

import (
	"context"

	"github.com/b97tsk/rx/internal/norec"
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

// ConcatWith applies Concat to the source Observable along with some other
// Observables to create a first-order Observable, then mirrors the resulting
// Observable.
func ConcatWith[T any](some ...Observable[T]) Operator[T, T] {
	return AsOperator(
		func(source Observable[T]) Observable[T] {
			return observables[T](append([]Observable[T]{source}, some...)).Concat
		},
	)
}

func (some observables[T]) Concat(ctx context.Context, sink Observer[T]) {
	var observer Observer[T]

	subscribeToNext := norec.Wrap(func() {
		if len(some) == 0 {
			sink.Complete()
			return
		}

		if err := ctx.Err(); err != nil {
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
