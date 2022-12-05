package rx

import (
	"context"
)

// Single emits the single item emitted by the source Observable.
// If the source emits more than one item or no items, Single throws
// ErrNotSingle or ErrEmpty respectively.
func Single[T any]() Operator[T, T] {
	return NewOperator(single[T])
}

func single[T any](source Observable[T]) Observable[T] {
	return func(ctx context.Context, sink Observer[T]) {
		ctx, cancel := context.WithCancel(ctx)

		sink = sink.WithCancel(cancel)

		var (
			value     T
			haveValue bool
		)

		var noop bool

		source.Subscribe(ctx, func(n Notification[T]) {
			switch {
			case noop:
			case n.HasValue:
				if !haveValue {
					value = n.Value
					haveValue = true
				} else {
					noop = true

					sink.Error(ErrNotSingle)
				}
			case n.HasError:
				sink(n)
			default:
				if haveValue {
					sink.Next(value)
					sink.Complete()
				} else {
					sink.Error(ErrEmpty)
				}
			}
		})
	}
}
