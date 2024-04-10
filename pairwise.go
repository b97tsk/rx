package rx

// Pairwise groups pairs of consecutive emissions together and emits them
// as Pairs.
func Pairwise[T any]() Operator[T, Pair[T, T]] {
	return NewOperator(
		func(source Observable[T]) Observable[Pair[T, T]] {
			return func(c Context, o Observer[Pair[T, T]]) {
				var p struct {
					Value    T
					HasValue bool
				}

				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext:
						if p.HasValue {
							o.Next(NewPair(p.Value, n.Value))
						}

						p.Value = n.Value
						p.HasValue = true

					case KindError:
						o.Error(n.Error)

					case KindComplete:
						o.Complete()
					}
				})
			}
		},
	)
}
