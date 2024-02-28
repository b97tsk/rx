package rx

// ToSlice collects all the values emitted by the source Observable,
// and then emits them as a slice when the source completes.
func ToSlice[T any]() Operator[T, []T] {
	return NewOperator(
		func(source Observable[T]) Observable[[]T] {
			return func(c Context, sink Observer[[]T]) {
				var s []T

				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext:
						s = append(s, n.Value)
					case KindError:
						sink.Error(n.Error)
					case KindComplete:
						Try1(sink, Next(s), func() { sink.Error(ErrOops) })
						sink.Complete()
					}
				})
			}
		},
	)
}
