package rx

// SkipWhile skips all values emitted by the source Observable as long as
// a given condition holds true, but emits all further source values as
// soon as the condition becomes false.
func SkipWhile[T any](cond func(v T) bool) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, sink Observer[T]) {
				var taking bool

				source.Subscribe(c, func(n Notification[T]) {
					switch {
					case taking || n.Kind != KindNext:
						sink(n)
					case !cond(n.Value):
						taking = true
						sink(n)
					}
				})
			}
		},
	)
}
