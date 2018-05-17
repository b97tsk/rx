package rx

import (
	"context"
)

type mapOperator struct {
	source  Operator
	project func(interface{}, int) interface{}
}

func (op mapOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	outerIndex := -1
	return op.source.Call(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			outerIndex++

			val := op.project(t.Value, outerIndex)
			ob.Next(val)

		case t.HasError:
			ob.Error(t.Value.(error))

		default:
			ob.Complete()
		}
	})
}

// Map creates an Observable that applies a given project function to each
// value emitted by the source Observable, then emits the resulting values.
func (o Observable) Map(project func(interface{}, int) interface{}) Observable {
	op := mapOperator{
		source:  o.Op,
		project: project,
	}
	return Observable{op}
}
