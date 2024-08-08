package rx

// DistinctComparable emits all values emitted by the source [Observable] that
// are distinct from each other.
func DistinctComparable[T comparable]() Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, o Observer[T]) {
				seen := make(map[T]struct{})

				source.Subscribe(c, func(n Notification[T]) {
					if n.Kind == KindNext {
						v := n.Value

						if _, exists := seen[v]; exists {
							return
						}

						seen[v] = struct{}{}
					}

					o.Emit(n)
				})
			}
		},
	)
}

// Distinct emits all values emitted by the source [Observable] whose mappings
// are distinct from each other.
func Distinct[T any, K comparable](mapping func(v T) K) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, o Observer[T]) {
				seen := make(map[K]struct{})

				source.Subscribe(c, func(n Notification[T]) {
					if n.Kind == KindNext {
						v := mapping(n.Value)

						if _, exists := seen[v]; exists {
							return
						}

						seen[v] = struct{}{}
					}

					o.Emit(n)
				})
			}
		},
	)
}
