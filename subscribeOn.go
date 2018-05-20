package rx

import (
	"context"
	"time"
)

type subscribeOnOperator struct {
	duration time.Duration
}

func (op subscribeOnOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	scheduleOnce(ctx, op.duration, func() {
		source.Subscribe(ctx, Finally(sink, cancel))
	})

	return ctx, cancel
}

// SubscribeOn creates an Observable that asynchronously subscribes the source
// Observable after waits for the duration to elapse.
func (o Observable) SubscribeOn(d time.Duration) Observable {
	op := subscribeOnOperator{d}
	return o.Lift(op.Call)
}
