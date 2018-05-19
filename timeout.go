package rx

import (
	"context"
	"time"
)

type timeoutOperator struct {
	timeout   time.Duration
	scheduler Scheduler
}

func (op timeoutOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	scheduleCancel := doNothing

	doSchedule := func() {
		scheduleCancel()

		_, scheduleCancel = op.scheduler.ScheduleOnce(ctx, op.timeout, func() {
			sink.Error(ErrTimeout)
			cancel()
		})
	}

	doSchedule()

	source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			sink(t)
			doSchedule()
		default:
			sink(t)
			cancel()
		}
	})

	return ctx, cancel
}

// Timeout creates an Observable that mirrors the source Observable or notify
// of an ErrTimeout if the source does not emit a value in given time span.
func (o Observable) Timeout(timeout time.Duration) Observable {
	op := timeoutOperator{timeout, DefaultScheduler}
	return o.Lift(op.Call).Mutex()
}
