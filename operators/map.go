package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// Map applies a given project function to each value emitted by the source,
// then emits the resulting values.
func Map(project func(interface{}, int) interface{}) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return func(ctx context.Context, sink rx.Observer) {
			sourceIndex := -1

			source.Subscribe(ctx, func(t rx.Notification) {
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

// MapTo emits the given constant value on the output Observable every time
// the source emits a value.
//
// It's like Map, but it maps every source value to the same output value
// every time.
func MapTo(val interface{}) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return func(ctx context.Context, sink rx.Observer) {
			source.Subscribe(ctx, func(t rx.Notification) {
				switch {
				case t.HasValue:
					sink.Next(val)
				default:
					sink(t)
				}
			})
		}
	}
}
