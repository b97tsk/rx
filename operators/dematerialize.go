package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

func dematerialize(source rx.Observable) rx.Observable {
	return func(ctx context.Context, sink rx.Observer) {
		ctx, cancel := context.WithCancel(ctx)
		sink = sink.WithCancel(cancel)

		var observer rx.Observer
		observer = func(t rx.Notification) {
			switch {
			case t.HasValue:
				if t, ok := t.Value.(rx.Notification); ok {
					switch {
					case t.HasValue:
						sink(t)
					default:
						observer = rx.Noop
						sink(t)
					}
				} else {
					observer = rx.Noop
					sink.Error(rx.ErrNotNotification)
				}
			default:
				sink(t)
			}
		}

		source.Subscribe(ctx, observer.Sink)
	}
}

// Dematerialize converts an Observable of Notification objects into the
// emissions that they represent. It's the opposite of Materialize.
func Dematerialize() rx.Operator {
	return dematerialize
}
