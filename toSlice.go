package rx

import (
	"context"
)

// ToSlice creates an Observable that collects all the values the source emits,
// then emits them as a slice when the source completes.
func (Operators) ToSlice() Operator {
	return func(source Observable) Observable {
		return func(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
			var values []interface{}
			return source.Subscribe(ctx, func(t Notification) {
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
}
