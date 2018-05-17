package rx

import (
	"context"
	"time"
)

type throwOperator struct {
	err       error
	delay     time.Duration
	scheduler Scheduler
}

func (op throwOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	if op.scheduler != nil {
		return op.scheduler.ScheduleOnce(ctx, op.delay, func() {
			ob.Error(op.err)
		})
	}

	ob.Error(op.err)
	return canceledCtx, noopFunc
}

// Throw creates an Observable that emits no items to the Observer and
// immediately emits an Error notification.
func Throw(err error) Observable {
	op := throwOperator{err: err}
	return Observable{}.Lift(op.Call)
}

// ThrowOn creates an Observable that emits no items to the Observer and
// immediately emits an Error notification, on the specified Scheduler.
func ThrowOn(err error, s Scheduler, delay time.Duration) Observable {
	op := throwOperator{err, delay, s}
	return Observable{}.Lift(op.Call)
}
