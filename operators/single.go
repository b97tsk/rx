package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// Single creates an Observable that emits the single item emitted by the
// source Observable. If the source emits more than one item or no items,
// throws rx.ErrNotSingle or rx.ErrEmpty respectively.
func Single() rx.Operator {
	return single
}

func single(source rx.Observable) rx.Observable {
	return func(ctx context.Context, sink rx.Observer) {
		ctx, cancel := context.WithCancel(ctx)

		sink = sink.WithCancel(cancel)

		var observer rx.Observer

		var (
			value    interface{}
			hasValue bool
		)

		observer = func(t rx.Notification) {
			switch {
			case t.HasValue:
				if hasValue {
					observer = rx.Noop
					sink.Error(rx.ErrNotSingle)
				} else {
					value = t.Value
					hasValue = true
				}

			case t.HasError:
				sink(t)

			default:
				if hasValue {
					sink.Next(value)
					sink.Complete()
				} else {
					sink.Error(rx.ErrEmpty)
				}
			}
		}

		source.Subscribe(ctx, observer.Sink)
	}
}
