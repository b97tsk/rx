package rx

// Reduce applies an accumulator function over the source [Observable],
// and emits the accumulated result when the source completes, given
// an initial value.
func Reduce[T, R any](init R, accumulator func(v1 R, v2 T) R) Operator[T, R] {
	return NewOperator(
		func(source Observable[T]) Observable[R] {
			return func(c Context, o Observer[R]) {
				v := init
				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext:
						v = accumulator(v, n.Value)
					case KindComplete:
						Try1(o, Next(v), func() { o.Stop(ErrOops) })
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
