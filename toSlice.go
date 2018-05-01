package rx

import (
	"context"
)

type toSliceOperator struct {
	source Operator
}

func (op toSliceOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	values := []interface{}(nil)
	return op.source.Call(ctx, ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			values = append(values, t.Value)
		case t.HasError:
			ob.Error(t.Value.(error))
		default:
			ob.Next(values)
			ob.Complete()
		}
	}))
}

// ToSlice creates an Observable that collects all the values the source emits,
// then emits them as a slice when the source completes.
func (o Observable) ToSlice() Observable {
	op := toSliceOperator{o.Op}
	return Observable{op}
}
