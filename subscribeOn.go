package rx

import (
	"context"
	"time"
)

// SubscribeOn creates an Observable that asynchronously subscribes to the
// source Observable after waits for the duration to elapse.
func (Operators) SubscribeOn(d time.Duration) Operator {
	return func(source Observable) Observable {
		return Create(
			func(ctx context.Context, sink Observer) {
				scheduleOnce(ctx, d, func() {
					source.Subscribe(ctx, sink)
				})
			},
		)
	}
}
