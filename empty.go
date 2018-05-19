package rx

import (
	"context"
	"time"
)

type emptyOperator struct {
	delay     time.Duration
	scheduler Scheduler
}

func (op emptyOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	if op.scheduler != nil {
		return op.scheduler.ScheduleOnce(ctx, op.delay, sink.Complete)
	}

	sink.Complete()
	return canceledCtx, doNothing
}

// Empty creates an Observable that emits no items to the Observer and
// immediately emits a Complete notification.
func Empty() Observable {
	op := emptyOperator{}
	return Observable{}.Lift(op.Call)
}

// EmptyOn creates an Observable that emits no items to the Observer and
// immediately emits a Complete notification, on the specified Scheduler.
func EmptyOn(s Scheduler, delay time.Duration) Observable {
	op := emptyOperator{delay, s}
	return Observable{}.Lift(op.Call)
}
