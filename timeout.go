package rx

import (
	"context"
	"time"
)

type timeoutOperator struct {
	source    Operator
	timeout   time.Duration
	scheduler Scheduler
}

func (op timeoutOperator) ApplyOptions(options []Option) Operator {
	for _, opt := range options {
		switch t := opt.(type) {
		case schedulerOption:
			op.scheduler = t.Value
		default:
			panic(ErrUnsupportedOption)
		}
	}
	return op
}

func (op timeoutOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	scheduleCancel := noopFunc

	doSchedule := func() {
		scheduleCancel()

		_, scheduleCancel = op.scheduler.ScheduleOnce(ctx, op.timeout, func() {
			ob.Error(ErrTimeout)
			cancel()
		})
	}

	doSchedule()

	op.source.Call(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			ob.Next(t.Value)
			doSchedule()
		case t.HasError:
			ob.Error(t.Value.(error))
			cancel()
		default:
			ob.Complete()
			cancel()
		}
	})

	return ctx, cancel
}

// Timeout creates an Observable that mirrors the source Observable or notify
// of an ErrTimeout if the source does not emit a value in given time span.
func (o Observable) Timeout(timeout time.Duration) Observable {
	op := timeoutOperator{
		source:    o.Op,
		timeout:   timeout,
		scheduler: DefaultScheduler,
	}
	return Observable{op}.Mutex()
}
