package rx

import (
	"context"
	"time"
)

type subscribeOnOperator struct {
	source    Operator
	delay     time.Duration
	scheduler Scheduler
}

func (op subscribeOnOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	op.scheduler.ScheduleOnce(ctx, op.delay, func() {
		op.source.Call(ctx, withFinalizer(ob, cancel))
	})

	return ctx, cancel
}

// SubscribeOn creates an Observable that asynchronously subscribes Observers
// to this Observable on the specified Scheduler.
func (o Observable) SubscribeOn(s Scheduler, delay time.Duration) Observable {
	op := subscribeOnOperator{
		source:    o.Op,
		delay:     delay,
		scheduler: s,
	}
	return Observable{op}
}
