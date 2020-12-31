package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// ToSlice collects all the values the source emits, and then emits them as
// a slice when the source completes.
func ToSlice() rx.Operator {
	return toSlice
}

func toSlice(source rx.Observable) rx.Observable {
	return func(ctx context.Context, sink rx.Observer) {
		var slice []interface{}

		source.Subscribe(ctx, func(t rx.Notification) {
			switch {
			case t.HasValue:
				slice = append(slice, t.Value)
			case t.HasError:
				sink(t)
			default:
				sink.Next(slice)
				sink.Complete()
			}
		})
	}
}
