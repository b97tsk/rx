package rx

import (
	"context"
)

// Last creates an Observable that emits only the last item emitted by the
// source Observable.
func (Operators) Last() Operator {
	return func(source Observable) Observable {
		return func(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
			var (
				lastValue    interface{}
				hasLastValue bool
			)
			return source.Subscribe(ctx, func(t Notification) {
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
						sink.Error(ErrEmpty)
					}
				}
			})
		}
	}
}
