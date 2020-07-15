package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

func isEmpty(source rx.Observable) rx.Observable {
	return func(ctx context.Context, sink rx.Observer) {
		ctx, cancel := context.WithCancel(ctx)
		sink = sink.WithCancel(cancel)

		var observer rx.Observer
		observer = func(t rx.Notification) {
			switch {
			case t.HasValue:
				observer = rx.Noop
				sink.Next(false)
				sink.Complete()
			case t.HasError:
				sink(t)
			default:
				sink.Next(true)
				sink.Complete()
			}
		}

		source.Subscribe(ctx, observer.Sink)
	}
}

// IsEmpty creates an Observable that emits true if the source Observable
// emits no items, otherwise, it emits false.
func IsEmpty() rx.Operator {
	return isEmpty
}
