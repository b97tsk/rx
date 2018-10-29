package rx

import (
	"context"
)

type toSliceOperator struct{}

func (op toSliceOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
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

// ToSlice creates an Observable that collects all the values the source emits,
// then emits them as a slice when the source completes.
func (Operators) ToSlice() OperatorFunc {
	return func(source Observable) Observable {
		op := toSliceOperator{}
		return source.Lift(op.Call)
	}
}
