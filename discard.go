package rx

// Discard ignores all values emitted by the source Observable.
func Discard[T any]() Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, sink Observer[T]) {
				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindError, KindComplete:
						sink(n)
					}
				})
			}
		},
	)
}

// IgnoreElements ignores all values emitted by the source Observable.
//
// It's like [Discard], but it can also change the output Observable to be
// of another type.
func IgnoreElements[T, R any]() Operator[T, R] {
	return NewOperator(
		func(source Observable[T]) Observable[R] {
			return func(c Context, sink Observer[R]) {
				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext:
					case KindError:
						sink.Error(n.Error)
					case KindComplete:
						sink.Complete()
					}
				})
			}
		},
	)
}