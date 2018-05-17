package rx

import (
	"context"
	"time"
)

type subscribeOnOperator struct {
	delay     time.Duration
	scheduler Scheduler
}

func (op subscribeOnOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	op.scheduler.ScheduleOnce(ctx, op.delay, func() {
		source.Subscribe(ctx, withFinalizer(ob, cancel))
	})

	return ctx, cancel
}

// SubscribeOn creates an Observable that asynchronously subscribes Observers
// to this Observable on the specified Scheduler.
func (o Observable) SubscribeOn(s Scheduler, delay time.Duration) Observable {
	op := subscribeOnOperator{delay, s}
	return o.Lift(op.Call)
}
