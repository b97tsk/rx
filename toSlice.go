package rx

import (
	"context"
)

type toSliceOperator struct{}

func (op toSliceOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	var values []interface{}
	return source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			values = append(values, t.Value)
		case t.HasError:
			t.Observe(ob)
		default:
			ob.Next(values)
			ob.Complete()
		}
	})
}

// ToSlice creates an Observable that collects all the values the source emits,
// then emits them as a slice when the source completes.
func (o Observable) ToSlice() Observable {
	op := toSliceOperator{}
	return o.Lift(op.Call)
}
