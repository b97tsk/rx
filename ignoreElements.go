package rx

import (
	"context"
)

// IgnoreElements creates an Observable that ignores all values emitted by the
// source Observable and only passes ERROR or COMPLETE emission.
func (Operators) IgnoreElements() Operator {
	return func(source Observable) Observable {
		return func(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
			return source.Subscribe(ctx, func(t Notification) {
				if t.HasValue {
					return
				}
				sink(t)
			})
		}
	}
}
