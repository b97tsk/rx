package rx

import (
	"context"
)

// Single emits the single value emitted by the source Observable.
// If the source emits more than one value or no values, it emits
// an error notification of ErrNotSingle or ErrEmpty respectively.
func Single[T any]() Operator[T, T] {
	return NewOperator(single[T])
}

func single[T any](source Observable[T]) Observable[T] {
	return func(ctx context.Context, sink Observer[T]) {
		ctx, cancel := context.WithCancel(ctx)

		sink = sink.OnLastNotification(cancel)

		var first struct {
			Value    T
			HasValue bool
		}

		var noop bool

		source.Subscribe(ctx, func(n Notification[T]) {
			switch {
			case noop:
			case n.HasValue:
				if !first.HasValue {
					first.Value = n.Value
					first.HasValue = true

					return
				}

				noop = true

				sink.Error(ErrNotSingle)

			case n.HasError:
				sink(n)

			default:
				if first.HasValue {
					sink.Next(first.Value)
					sink.Complete()
				} else {
					sink.Error(ErrEmpty)
				}
			}
		})
	}
}
