package rx

import (
	"context"
)

// Map creates an Observable that applies a given project function to each
// value emitted by the source Observable, then emits the resulting values.
func (Operators) Map(project func(interface{}, int) interface{}) Operator {
	return func(source Observable) Observable {
		return func(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
			var sourceIndex = -1
			return source.Subscribe(ctx, func(t Notification) {
				switch {
				case t.HasValue:
					sourceIndex++
					sink.Next(project(t.Value, sourceIndex))
				default:
					sink(t)
				}
			})
		}
	}
}

// MapTo creates an Observable that emits the given constant value on the
// output Observable every time the source Observable emits a value.
//
// It's like Map, but it maps every source value to the same output value
// every time.
func (Operators) MapTo(value interface{}) Operator {
	return func(source Observable) Observable {
		return func(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
			return source.Subscribe(ctx, func(t Notification) {
				switch {
				case t.HasValue:
					sink.Next(value)
				default:
					sink(t)
				}
			})
		}
	}
}
