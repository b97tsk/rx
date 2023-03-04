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

// OnNext mirrors the source Observable, passing values to f before
// each value emission.
func OnNext[T any](f func(v T)) Operator[T, T] {
	if f == nil {
		panic("f == nil")
	}

	return onNext(f)
}

func onNext[T any](f func(v T)) Operator[T, T] {
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

// OnComplete mirrors the source Observable, and calls f when the source
// completes.
func OnComplete[T any](f func()) Operator[T, T] {
	if f == nil {
		panic("f == nil")
	}

	return onComplete[T](f)
}

func onComplete[T any](f func()) Operator[T, T] {
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

// OnError mirrors the source Observable, and calls f when the source emits
// a notification of error.
func OnError[T any](f func(err error)) Operator[T, T] {
	if f == nil {
		panic("f == nil")
	}

	return onError[T](f)
}

func onError[T any](f func(err error)) Operator[T, T] {
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

// OnLastNotification mirrors the source Observable, and calls f when
// the source completes or emits a notification of error.
func OnLastNotification[T any](f func()) Operator[T, T] {
	if f == nil {
		panic("f == nil")
	}

	return onLastNotification[T](f)
}

func onLastNotification[T any](f func()) Operator[T, T] {
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
