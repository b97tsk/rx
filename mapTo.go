package rx

import (
	"context"
)

type mapToOperator struct {
	value interface{}
}

func (op mapToOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	return source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			ob.Next(op.value)
		default:
			t.Observe(ob)
		}
	})
}

// MapTo creates an Observable that emits the given constant value on the
// output Observable every time the source Observable emits a value.
//
// It's like Map, but it maps every source value to the same output value
// every time.
func (o Observable) MapTo(value interface{}) Observable {
	op := mapToOperator{value}
	return o.Lift(op.Call)
}
