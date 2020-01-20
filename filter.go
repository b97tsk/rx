package rx

import (
	"context"
)

// Filter creates an Observable that filter items emitted by the source
// Observable by only emitting those that satisfy a specified predicate.
func (Operators) Filter(predicate func(interface{}, int) bool) Operator {
	return func(source Observable) Observable {
		return func(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
			var sourceIndex = -1
			return source.Subscribe(ctx, func(t Notification) {
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

// Exclude creates an Observable that filter items emitted by the source
// Observable by only emitting those that do not satisfy a specified predicate.
func (Operators) Exclude(predicate func(interface{}, int) bool) Operator {
	return func(source Observable) Observable {
		return func(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
			var sourceIndex = -1
			return source.Subscribe(ctx, func(t Notification) {
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
