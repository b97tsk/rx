package rx

import (
	"context"
)

type mapObservable struct {
	Source  Observable
	Project func(interface{}, int) interface{}
}

func (obs mapObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	var outerIndex = -1
	return obs.Source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			outerIndex++

			val := obs.Project(t.Value, outerIndex)
			sink.Next(val)

		default:
			sink(t)
		}
	})
}

// Map creates an Observable that applies a given project function to each
// value emitted by the source Observable, then emits the resulting values.
func (Operators) Map(project func(interface{}, int) interface{}) Operator {
	return func(source Observable) Observable {
		return mapObservable{source, project}.Subscribe
	}
}
