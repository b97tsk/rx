package rx

// CompactComparable emits all values emitted by the source [Observable] that
// are distinct from the previous.
func CompactComparable[T comparable]() Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, o Observer[T]) {
				var last struct {
					value    T
					hasValue bool
				}

				source.Subscribe(c, func(n Notification[T]) {
					if n.Kind == KindNext {
						if last.hasValue && last.value == n.Value {
							return
						}

						last.value = n.Value
						last.hasValue = true
					}

					o.Emit(n)
				})
			}
		},
	)
}

// Compact emits all values emitted by the source [Observable] that are
// distinct from the previous, given a comparison function.
func Compact[T any](eq func(v1, v2 T) bool) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, o Observer[T]) {
				var last struct {
					value    T
					hasValue bool
				}

				source.Subscribe(c, func(n Notification[T]) {
					if n.Kind == KindNext {
						if last.hasValue && eq(last.value, n.Value) {
							return
						}

						last.value = n.Value
						last.hasValue = true
					}

					o.Emit(n)
				})
			}
		},
	)
}

// CompactComparableKey emits all values emitted by the source [Observable]
// whose mappings are distinct from the previous.
func CompactComparableKey[T any, K comparable](mapping func(v T) K) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, o Observer[T]) {
				var last struct {
					value    K
					hasValue bool
				}

				source.Subscribe(c, func(n Notification[T]) {
					if n.Kind == KindNext {
						keyValue := mapping(n.Value)

						if last.hasValue && last.value == keyValue {
							return
						}

						last.value = keyValue
						last.hasValue = true
					}

					o.Emit(n)
				})
			}
		},
	)
}

// CompactKey emits all values emitted by the source [Observable] whose
// mappings are distinct from the previous, given a comparison function.
func CompactKey[T, K any](mapping func(v T) K, eq func(v1, v2 K) bool) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, o Observer[T]) {
				var last struct {
					value    K
					hasValue bool
				}

				source.Subscribe(c, func(n Notification[T]) {
					if n.Kind == KindNext {
						keyValue := mapping(n.Value)

						if last.hasValue && eq(last.value, keyValue) {
							return
						}

						last.value = keyValue
						last.hasValue = true
					}

					o.Emit(n)
				})
			}
		},
	)
}
