package rx

import (
	"context"
)

type toSliceOperator struct {
	source Operator
}

func (op toSliceOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	all := []interface{}(nil)
	return op.source.Call(ctx, ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			all = append(all, t.Value)
		case t.HasError:
			ob.Error(t.Value.(error))
		default:
			ob.Next(all)
			ob.Complete()
		}
	}))
}

// ToSlice creates an Observable that collects all source emissions and emits
// them as a slice when the source completes.
func (o Observable) ToSlice() Observable {
	op := toSliceOperator{o.Op}
	return Observable{op}
}
