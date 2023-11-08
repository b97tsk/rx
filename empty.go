package rx

import "context"

// Empty returns an Observable that emits no values and immediately completes.
func Empty[T any]() Observable[T] {
	return empty[T]
}

func empty[T any](_ context.Context, sink Observer[T]) {
	sink.Complete()
}

// Never returns an Observable that never emits anything, except
// when a context cancellation is detected, emits an error notification
// of whatever that context reports.
func Never[T any]() Observable[T] {
	return never[T]
}

func never[T any](ctx context.Context, sink Observer[T]) {
	if ctx.Done() != nil {
		wg := WaitGroupFromContext(ctx)
		if wg != nil {
			wg.Add(1)
		}

		context.AfterFunc(ctx, func() {
			if wg != nil {
				defer wg.Done()
			}

			sink.Error(ctx.Err())
		})
	}
}

// Throw creates an Observable that emits no values and immediately emits
// an error notification of err.
func Throw[T any](err error) Observable[T] {
	return func(ctx context.Context, sink Observer[T]) {
		sink.Error(err)
	}
}
