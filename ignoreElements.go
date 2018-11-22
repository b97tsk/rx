package rx

import (
	"context"
)

// IgnoreElements creates an Observable that ignores all values emitted by the
// source Observable and only passes Complete or Error emission.
func (Operators) IgnoreElements() OperatorFunc {
	return func(source Observable) Observable {
		return source.Lift(
			func(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
				return source.Subscribe(ctx, func(t Notification) {
					if t.HasValue {
						return
					}
					sink(t)
				})
			},
		)
	}
}
