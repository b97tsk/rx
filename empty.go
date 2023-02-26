package rx

import (
	"context"

	"github.com/b97tsk/rx/internal/ctxwatch"
	"github.com/b97tsk/rx/internal/waitgroup"
)

// Empty returns an Observable that emits no items to the Observer and
// immediately completes.
func Empty[T any]() Observable[T] {
	return empty[T]
}

func empty[T any](_ context.Context, sink Observer[T]) {
	sink.Complete()
}

// Never returns an Observable that never emits anything.
func Never[T any]() Observable[T] {
	return never[T]
}

func never[T any](ctx context.Context, sink Observer[T]) {
	if ctx.Done() != nil {
		wg := waitgroup.Get(ctx)
		if wg != nil {
			wg.Add(1)
		}

		ctxwatch.Add(ctx, func() {
			if wg != nil {
				defer wg.Done()
			}

			sink.Error(ctx.Err())
		})
	}
}

// Throw creates an Observable that emits no items to the Observer and
// immediately throws a specified error.
func Throw[T any](err error) Observable[T] {
	return func(ctx context.Context, sink Observer[T]) {
		sink.Error(err)
	}
}
