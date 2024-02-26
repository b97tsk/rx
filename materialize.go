package rx

// Materialize represents all of the Notifications from the source Observable
// as values, and then completes.
func Materialize[T any]() Operator[T, Notification[T]] {
	return NewOperator(materialize[T])
}

func materialize[T any](source Observable[T]) Observable[Notification[T]] {
	return func(c Context, sink Observer[Notification[T]]) {
		source.Subscribe(c, func(n Notification[T]) {
			sink.Next(n)
			switch n.Kind {
			case KindError, KindComplete:
				sink.Complete()
			}
		})
	}
}
