package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// Filter creates an Observable that filter items emitted by the source
// Observable by only emitting those that satisfy a specified predicate.
func Filter(predicate func(interface{}, int) bool) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return func(ctx context.Context, sink rx.Observer) {
			sourceIndex := -1
			source.Subscribe(ctx, func(t rx.Notification) {
				switch {
				case t.HasValue:
					sourceIndex++
					if predicate(t.Value, sourceIndex) {
						sink(t)
					}
				default:
					sink(t)
				}
			})
		}
	}
}

// FilterMap creates an Observable that passes each item emitted by the source
// Observable to a specified predicate and emits their mapping, the first
// return value of the predicate, if the second is true.
func FilterMap(predicate func(interface{}, int) (interface{}, bool)) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return func(ctx context.Context, sink rx.Observer) {
			sourceIndex := -1
			source.Subscribe(ctx, func(t rx.Notification) {
				switch {
				case t.HasValue:
					sourceIndex++
					if val, ok := predicate(t.Value, sourceIndex); ok {
						sink.Next(val)
					}
				default:
					sink(t)
				}
			})
		}
	}
}

// Exclude creates an Observable that filter items emitted by the source
// Observable by only emitting those that do not satisfy a specified predicate.
func Exclude(predicate func(interface{}, int) bool) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return func(ctx context.Context, sink rx.Observer) {
			sourceIndex := -1
			source.Subscribe(ctx, func(t rx.Notification) {
				switch {
				case t.HasValue:
					sourceIndex++
					if !predicate(t.Value, sourceIndex) {
						sink(t)
					}
				default:
					sink(t)
				}
			})
		}
	}
}
