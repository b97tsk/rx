package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// First emits only the first value emitted by the source, if the source
// turns out to be empty, throws ErrEmpty.
func First() rx.Operator {
	return first
}

func first(source rx.Observable) rx.Observable {
	return func(ctx context.Context, sink rx.Observer) {
		ctx, cancel := context.WithCancel(ctx)

		sink = sink.WithCancel(cancel)

		var observer rx.Observer

		observer = func(t rx.Notification) {
			switch {
			case t.HasValue:
				observer = rx.Noop

				sink(t)
				sink.Complete()

			case t.HasError:
				sink(t)

			default:
				sink.Error(rx.ErrEmpty)
			}
		}

		source.Subscribe(ctx, observer.Sink)
	}
}

// FirstOrDefault emits only the first value emitted by the source, if the
// source turns out to be empty, emits a specified default value.
func FirstOrDefault(def interface{}) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return func(ctx context.Context, sink rx.Observer) {
			ctx, cancel := context.WithCancel(ctx)

			sink = sink.WithCancel(cancel)

			var observer rx.Observer

			observer = func(t rx.Notification) {
				switch {
				case t.HasValue:
					observer = rx.Noop

					sink(t)
					sink.Complete()

				case t.HasError:
					sink(t)

				default:
					sink.Next(def)
					sink(t)
				}
			}

			source.Subscribe(ctx, observer.Sink)
		}
	}
}
