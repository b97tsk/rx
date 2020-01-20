package rx

import (
	"context"
)

// Count creates an Observable that counts the number of NEXT emissions on
// the source and emits that number when the source completes.
func (Operators) Count() Operator {
	return func(source Observable) Observable {
		return func(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
			var count int
			return source.Subscribe(ctx, func(t Notification) {
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
}
