package rx

import (
	"context"
)

type mapToOperator struct {
	source Operator
	value  interface{}
}

func (op mapToOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	return op.source.Call(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			ob.Next(op.value)
		case t.HasError:
			ob.Error(t.Value.(error))
		default:
			ob.Complete()
		}
	})
}

// MapTo creates an Observable that emits the given constant value on the
// output Observable every time the source Observable emits a value.
//
// It's like Map, but it maps every source value to the same output value
// every time.
func (o Observable) MapTo(value interface{}) Observable {
	op := mapToOperator{
		source: o.Op,
		value:  value,
	}
	return Observable{op}
}
