package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

func first(source rx.Observable) rx.Observable {
	return rx.Create(
		func(ctx context.Context, sink rx.Observer) {
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
			source.Subscribe(ctx, observer.Notify)
		},
	)
}

// First creates an Observable that emits only the first value (or the first
// value that meets some condition) emitted by the source Observable.
func First() rx.Operator {
	return first
}
