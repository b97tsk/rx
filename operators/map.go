package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// Map creates an Observable that applies a given project function to each
// value emitted by the source Observable, then emits the resulting values.
func Map(project func(interface{}, int) interface{}) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return func(ctx context.Context, sink rx.Observer) (context.Context, context.CancelFunc) {
			sourceIndex := -1
			return source.Subscribe(ctx, func(t rx.Notification) {
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
func MapTo(value interface{}) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return func(ctx context.Context, sink rx.Observer) (context.Context, context.CancelFunc) {
			return source.Subscribe(ctx, func(t rx.Notification) {
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
