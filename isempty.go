package rx

// IsEmpty emits a boolean to indicate whether the source emits no values.
func IsEmpty[T any]() Operator[T, bool] {
	return NewOperator(isEmpty[T])
}

func isEmpty[T any](source Observable[T]) Observable[bool] {
	return func(c Context, sink Observer[bool]) {
		c, cancel := c.WithCancel()

		var noop bool

		source.Subscribe(c, func(n Notification[T]) {
			if noop {
				return
			}

			switch n.Kind {
			case KindNext, KindError, KindComplete:
				noop = true
				cancel()

				switch n.Kind {
				case KindNext:
					sink.Next(false)
					sink.Complete()
				case KindError:
					sink.Error(n.Error)
				case KindComplete:
					sink.Next(true)
					sink.Complete()
				}
			}
		})
	}
}
