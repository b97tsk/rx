package rx

import (
	"context"
)

// Dematerialize converts an Observable of Notification values into
// the emissions that they represent. It's the opposite of [Materialize].
func Dematerialize[_ Notification[T], T any]() Operator[Notification[T], T] {
	return NewOperator(dematerialize[Notification[T]])
}

func dematerialize[_ Notification[T], T any](source Observable[Notification[T]]) Observable[T] {
	return func(ctx context.Context, sink Observer[T]) {
		ctx, cancel := context.WithCancel(ctx)

		sink = sink.OnLastNotification(cancel)

		var noop bool

		source.Subscribe(ctx, func(n Notification[Notification[T]]) {
			if noop {
				return
			}

			switch n.Kind {
			case KindNext:
				n := n.Value

				switch n.Kind {
				case KindError, KindComplete:
					noop = true
				}

				sink(n)
			case KindError:
				sink.Error(n.Error)
			case KindComplete:
				sink.Complete()
			}
		})
	}
}
