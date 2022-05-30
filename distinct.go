package rx

import (
	"context"
)

// DistinctComparable emits all items emitted by the source Observable that
// are distinct from each other.
func DistinctComparable[T comparable]() Operator[T, T] {
	return AsOperator(
		func(source Observable[T]) Observable[T] {
			return func(ctx context.Context, sink Observer[T]) {
				seen := make(map[T]struct{})

				source.Subscribe(ctx, func(n Notification[T]) {
					if n.HasValue {
						v := n.Value

						if _, exists := seen[v]; exists {
							return
						}

						seen[v] = struct{}{}
					}

					sink(n)
				})
			}
		},
	)
}

// Distinct emits all items emitted by the source Observable whose projections
// are distinct from each other.
func Distinct[T any, K comparable](proj func(v T) K) Operator[T, T] {
	if proj == nil {
		panic("proj == nil")
	}

	return distinct(proj)
}

func distinct[T any, K comparable](proj func(v T) K) Operator[T, T] {
	return AsOperator(
		func(source Observable[T]) Observable[T] {
			return func(ctx context.Context, sink Observer[T]) {
				seen := make(map[K]struct{})

				source.Subscribe(ctx, func(n Notification[T]) {
					if n.HasValue {
						v := proj(n.Value)

						if _, exists := seen[v]; exists {
							return
						}

						seen[v] = struct{}{}
					}

					sink(n)
				})
			}
		},
	)
}
