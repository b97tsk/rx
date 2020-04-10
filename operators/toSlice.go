package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

func toSlice(source rx.Observable) rx.Observable {
	return func(ctx context.Context, sink rx.Observer) (context.Context, context.CancelFunc) {
		var values []interface{}
		return source.Subscribe(ctx, func(t rx.Notification) {
			switch {
			case t.HasValue:
				values = append(values, t.Value)
			case t.HasError:
				sink(t)
			default:
				sink.Next(values)
				sink.Complete()
			}
		})
	}
}

// ToSlice creates an Observable that collects all the values the source emits,
// then emits them as a slice when the source completes.
func ToSlice() rx.Operator {
	return toSlice
}
