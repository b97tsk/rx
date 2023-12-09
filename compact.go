package rx

import (
	"context"
)

// CompactComparable emits all values emitted by the source Observable that
// are distinct from the previous.
func CompactComparable[T comparable]() Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(ctx context.Context, sink Observer[T]) {
				var last struct {
					Value    T
					HasValue bool
				}

				source.Subscribe(ctx, func(n Notification[T]) {
					if n.Kind == KindNext {
						if last.HasValue && last.Value == n.Value {
							return
						}

						last.Value = n.Value
						last.HasValue = true
					}

					sink(n)
				})
			}
		},
	)
}

// Compact emits all values emitted by the source Observable that are distinct
// from the previous, given a comparison function.
func Compact[T any](eq func(v1, v2 T) bool) Operator[T, T] {
	if eq == nil {
		panic("eq == nil")
	}

	return compact(eq)
}

func compact[T any](eq func(v1, v2 T) bool) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(ctx context.Context, sink Observer[T]) {
				var last struct {
					Value    T
					HasValue bool
				}

				source.Subscribe(ctx, func(n Notification[T]) {
					if n.Kind == KindNext {
						if last.HasValue && eq(last.Value, n.Value) {
							return
						}

						last.Value = n.Value
						last.HasValue = true
					}

					sink(n)
				})
			}
		},
	)
}

// CompactComparableKey emits all values emitted by the source Observable whose
// projections are distinct from the previous.
func CompactComparableKey[T any, K comparable](proj func(v T) K) Operator[T, T] {
	if proj == nil {
		panic("proj == nil")
	}

	return compactComparableKey(proj)
}

func compactComparableKey[T any, K comparable](proj func(v T) K) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(ctx context.Context, sink Observer[T]) {
				var last struct {
					Value    K
					HasValue bool
				}

				source.Subscribe(ctx, func(n Notification[T]) {
					if n.Kind == KindNext {
						keyValue := proj(n.Value)

						if last.HasValue && last.Value == keyValue {
							return
						}

						last.Value = keyValue
						last.HasValue = true
					}

					sink(n)
				})
			}
		},
	)
}

// CompactKey emits all values emitted by the source Observable whose
// projections are distinct from the previous, given a comparison function.
func CompactKey[T, K any](proj func(v T) K, eq func(v1, v2 K) bool) Operator[T, T] {
	switch {
	case proj == nil:
		panic("proj == nil")
	case eq == nil:
		panic("eq == nil")
	}

	return compactKey(proj, eq)
}

func compactKey[T, K any](proj func(v T) K, eq func(v1, v2 K) bool) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(ctx context.Context, sink Observer[T]) {
				var last struct {
					Value    K
					HasValue bool
				}

				source.Subscribe(ctx, func(n Notification[T]) {
					if n.Kind == KindNext {
						keyValue := proj(n.Value)

						if last.HasValue && eq(last.Value, keyValue) {
							return
						}

						last.Value = keyValue
						last.HasValue = true
					}

					sink(n)
				})
			}
		},
	)
}
