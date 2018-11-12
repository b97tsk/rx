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

	sink = Finally(sink, cancel)

	var (
		scheduleCancel = nothingToDo

		try cancellableLocker
	)

	doNextAndUnlock := func(t Notification) {
		defer try.Unlock()
		sink(t)
	}

	doSchedule := func() {
		scheduleCancel()

		_, scheduleCancel = scheduleOnce(ctx, op.Duration, func() {
			if try.Lock() {
				try.CancelAndUnlock()
				sink.Error(ErrTimeout)
			}
		})
	}

	doSchedule()

	source.Subscribe(ctx, func(t Notification) {
		if try.Lock() {
			switch {
			case t.HasValue:
				doNextAndUnlock(t)
				doSchedule()
			default:
				try.CancelAndUnlock()
				sink(t)
			}
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
