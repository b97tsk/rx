package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// Count creates an Observable that counts the number of values on the source
// and emits that number when the source completes.
func Count() rx.Operator {
	return count
}

func count(source rx.Observable) rx.Observable {
	return func(ctx context.Context, sink rx.Observer) {
		var count int
		source.Subscribe(ctx, func(t rx.Notification) {
			switch {
			case t.HasValue:
				count++
			case t.HasError:
				sink(t)
			default:
				sink.Next(count)
				sink.Complete()
			}
		})
	}
}
