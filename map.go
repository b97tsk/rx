package rx

import (
	"context"
)

type mapOperator struct {
	Project func(interface{}, int) interface{}
}

func (op mapOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	var outerIndex = -1
	return source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			outerIndex++

			val := op.Project(t.Value, outerIndex)
			sink.Next(val)

		default:
			sink(t)
		}
	})
}

// Map creates an Observable that applies a given project function to each
// value emitted by the source Observable, then emits the resulting values.
func (Operators) Map(project func(interface{}, int) interface{}) OperatorFunc {
	return func(source Observable) Observable {
		op := mapOperator{project}
		return source.Lift(op.Call)
	}
}
