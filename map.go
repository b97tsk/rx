package rx

import (
	"context"
)

type mapOperator struct {
	project func(interface{}, int) interface{}
}

func (op mapOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	var outerIndex = -1
	return source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			outerIndex++

			val := op.project(t.Value, outerIndex)
			ob.Next(val)

		default:
			t.Observe(ob)
		}
	})
}

// Map creates an Observable that applies a given project function to each
// value emitted by the source Observable, then emits the resulting values.
func (o Observable) Map(project func(interface{}, int) interface{}) Observable {
	op := mapOperator{project}
	return o.Lift(op.Call)
}
