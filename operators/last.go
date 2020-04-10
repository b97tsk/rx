package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

func last(source rx.Observable) rx.Observable {
	return func(ctx context.Context, sink rx.Observer) (context.Context, context.CancelFunc) {
		var (
			lastValue    interface{}
			hasLastValue bool
		)
		return source.Subscribe(ctx, func(t rx.Notification) {
			switch {
			case t.HasValue:
				lastValue = t.Value
				hasLastValue = true
			case t.HasError:
				sink(t)
			default:
				if hasLastValue {
					sink.Next(lastValue)
					sink.Complete()
				} else {
					sink.Error(rx.ErrEmpty)
				}
			}
		})
	}
}

// Last creates an Observable that emits only the last item emitted by the
// source Observable.
func Last() rx.Operator {
	return last
}
