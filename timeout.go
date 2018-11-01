package rx

import (
	"context"
	"time"
)

type timeoutOperator struct {
	Duration time.Duration
}

func (op timeoutOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	scheduleCancel := nothingToDo

	sink = Mutex(Finally(sink, cancel))

	doSchedule := func() {
		scheduleCancel()

		_, scheduleCancel = scheduleOnce(ctx, op.Duration, func() {
			sink.Error(ErrTimeout)
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
		}
	})

	return ctx, cancel
}

// Timeout creates an Observable that mirrors the source Observable or notify
// of an ErrTimeout if the source does not emit a value in given time span.
func (Operators) Timeout(timeout time.Duration) OperatorFunc {
	return func(source Observable) Observable {
		op := timeoutOperator{timeout}
		return source.Lift(op.Call)
	}
}
