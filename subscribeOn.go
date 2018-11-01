package rx

import (
	"context"
	"time"
)

type subscribeOnOperator struct {
	Duration time.Duration
}

func (op subscribeOnOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	scheduleOnce(ctx, op.Duration, func() {
		source.Subscribe(ctx, sink)
	})

	return ctx, cancel
}

// SubscribeOn creates an Observable that asynchronously subscribes the source
// Observable after waits for the duration to elapse.
func (Operators) SubscribeOn(d time.Duration) OperatorFunc {
	return func(source Observable) Observable {
		op := subscribeOnOperator{d}
		return source.Lift(op.Call)
	}
}
