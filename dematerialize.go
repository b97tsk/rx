package rx

import (
	"context"
)

// Dematerialize converts an Observable of Notification objects into the
// emissions that they represent. It's the opposite of Materialize.
func Dematerialize[_ Notification[T], T any]() Operator[Notification[T], T] {
	return AsOperator(dematerialize[Notification[T]])
}

func dematerialize[_ Notification[T], T any](source Observable[Notification[T]]) Observable[T] {
	return func(ctx context.Context, sink Observer[T]) {
		ctx, cancel := context.WithCancel(ctx)

		sink = sink.WithCancel(cancel)

		var noop bool

		source.Subscribe(ctx, func(n Notification[Notification[T]]) {
			switch {
			case noop:
			case n.HasValue:
				n := n.Value
				noop = !n.HasValue

				sink(n)
			case n.HasError:
				sink.Error(n.Error)
			default:
				sink.Complete()
			}
		})
	}
}
