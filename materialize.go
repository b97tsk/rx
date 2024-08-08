package rx

// Materialize represents all of the notifications, excluding [Stop]
// notifications, from the source [Observable] as values, and then completes.
func Materialize[T any]() Operator[T, Notification[T]] {
	return NewOperator(
		func(source Observable[T]) Observable[Notification[T]] {
			return func(c Context, o Observer[Notification[T]]) {
				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext:
						o.Next(n)
					case KindComplete, KindError:
						Try1(o, Next(n), func() { o.Stop(ErrOops) })
						o.Complete()
					case KindStop:
						o.Stop(n.Error)
					}
				})
			}
		},
	)
}
