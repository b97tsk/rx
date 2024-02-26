package rx

// Last emits only the last value emitted by the source Observable.
// If the source turns out to be empty, it emits an error notification
// of ErrEmpty.
func Last[T any]() Operator[T, T] {
	return NewOperator(last[T])
}

func last[T any](source Observable[T]) Observable[T] {
	return func(c Context, sink Observer[T]) {
		var last struct {
			Value    T
			HasValue bool
		}

		source.Subscribe(c, func(n Notification[T]) {
			switch n.Kind {
			case KindNext:
				last.Value = n.Value
				last.HasValue = true
			case KindError:
				sink(n)
			case KindComplete:
				if last.HasValue {
					sink.Next(last.Value)
					sink.Complete()
				} else {
					sink.Error(ErrEmpty)
				}
			}
		})
	}
}

// LastOrElse emits only the last value emitted by the source Observable.
// If the source turns out to be empty, it emits a specified default value.
func LastOrElse[T any](def T) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return func(c Context, sink Observer[T]) {
				var last struct {
					Value    T
					HasValue bool
				}

				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext:
						last.Value = n.Value
						last.HasValue = true
					case KindError:
						sink(n)
					case KindComplete:
						if last.HasValue {
							sink.Next(last.Value)
						} else {
							sink.Next(def)
						}

						sink(n)
					}
				})
			}
		},
	)
}
