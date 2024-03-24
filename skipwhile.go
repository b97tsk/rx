package rx

// SkipWhile skips all values emitted by the source Observable as long as
// a given predicate function returns true, but emits all further source
// values as soon as the predicate function returns false.
func SkipWhile[T any](pred func(v T) bool) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, sink Observer[T]) {
				var taking bool

				source.Subscribe(c, func(n Notification[T]) {
					switch {
					case taking || n.Kind != KindNext:
						sink(n)
					case !pred(n.Value):
						taking = true
						sink(n)
					}
				})
			}
		},
	)
}
