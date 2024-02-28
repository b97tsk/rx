package rx

// Flat flattens a higher-order Observable into a first-order Observable,
// by applying a flat function to the inner Observables.
func Flat[_ Observable[T], T any](f func(some ...Observable[T]) Observable[T]) Operator[Observable[T], T] {
	return NewOperator(
		func(source Observable[Observable[T]]) Observable[T] {
			return func(c Context, sink Observer[T]) {
				var s []Observable[T]

				source.Subscribe(c, func(n Notification[Observable[T]]) {
					switch n.Kind {
					case KindNext:
						s = append(s, n.Value)
					case KindError:
						sink.Error(n.Error)
					case KindComplete:
						Try01(
							func() Observable[T] { return f(s...) },
							func() { sink.Error(ErrOops) },
						).Subscribe(c, sink)
					}
				})
			}
		},
	)
}
