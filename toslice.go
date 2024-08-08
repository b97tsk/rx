package rx

// ToSlice collects all the values emitted by the source [Observable],
// and then emits them as a slice when the source completes.
func ToSlice[T any]() Operator[T, []T] {
	return NewOperator(
		func(source Observable[T]) Observable[[]T] {
			return func(c Context, o Observer[[]T]) {
				var s []T
				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext:
						s = append(s, n.Value)
					case KindComplete:
						Try1(o, Next(s), func() { o.Stop(ErrOops) })
						o.Complete()
					case KindError:
						o.Error(n.Error)
					case KindStop:
						o.Stop(n.Error)
					}
				})
			}
		},
	)
}
