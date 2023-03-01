package rx

import (
	"context"
)

// Do mirrors the source Observable, passing emissions to tap before
// each emission.
func Do[T any](tap Observer[T]) Operator[T, T] {
	if tap == nil {
		panic("tap == nil")
	}

	return do(tap)
}

func do[T any](tap Observer[T]) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(ctx context.Context, sink Observer[T]) {
				source.Subscribe(ctx, func(n Notification[T]) {
					tap(n)
					sink(n)
				})
			}
		},
	)
}

// DoOnNext mirrors the source Observable, passing values to f before
// each value emission.
func DoOnNext[T any](f func(v T)) Operator[T, T] {
	if f == nil {
		panic("f == nil")
	}

	return doOnNext(f)
}

func doOnNext[T any](f func(v T)) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(ctx context.Context, sink Observer[T]) {
				source.Subscribe(ctx, func(n Notification[T]) {
					if n.HasValue {
						f(n.Value)
					}

					sink(n)
				})
			}
		},
	)
}

// DoOnComplete mirrors the source Observable, and calls f when the source
// completes.
func DoOnComplete[T any](f func()) Operator[T, T] {
	if f == nil {
		panic("f == nil")
	}

	return doOnComplete[T](f)
}

func doOnComplete[T any](f func()) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(ctx context.Context, sink Observer[T]) {
				source.Subscribe(ctx, func(n Notification[T]) {
					if !n.HasValue && !n.HasError {
						f()
					}

					sink(n)
				})
			}
		},
	)
}

// DoOnError mirrors the source Observable, and calls f when the source emits
// a notification of error.
func DoOnError[T any](f func(err error)) Operator[T, T] {
	if f == nil {
		panic("f == nil")
	}

	return doOnError[T](f)
}

func doOnError[T any](f func(err error)) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(ctx context.Context, sink Observer[T]) {
				source.Subscribe(ctx, func(n Notification[T]) {
					if n.HasError {
						f(n.Error)
					}

					sink(n)
				})
			}
		},
	)
}

// DoOnErrorOrComplete mirrors the source Observable, and calls f when
// the source completes or emits a notification of error.
func DoOnErrorOrComplete[T any](f func()) Operator[T, T] {
	if f == nil {
		panic("f == nil")
	}

	return doOnErrorOrComplete[T](f)
}

func doOnErrorOrComplete[T any](f func()) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(ctx context.Context, sink Observer[T]) {
				source.Subscribe(ctx, func(n Notification[T]) {
					if !n.HasValue {
						f()
					}

					sink(n)
				})
			}
		},
	)
}
