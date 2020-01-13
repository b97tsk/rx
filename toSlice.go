package rx

import (
	"context"
)

type toSliceObservable struct {
	Source Observable
}

func (obs toSliceObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	var values []interface{}
	return obs.Source.Subscribe(ctx, func(t Notification) {
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
		return toSliceObservable{source}.Subscribe
	}
}
