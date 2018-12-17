package rx

import (
	"context"
	"time"
)

// A TimeoutConfigure is a configure for Timeout.
type TimeoutConfigure struct {
	Duration   time.Duration
	Observable Observable
}

// MakeFunc creates an OperatorFunc from this type.
func (conf TimeoutConfigure) MakeFunc() OperatorFunc {
	return MakeFunc(timeoutOperator(conf).Call)
}

type timeoutOperator TimeoutConfigure

func (op timeoutOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	var try cancellableLocker
	_, scheduleCancel := Done()
	childCtx, childCancel := context.WithCancel(ctx)

	doSchedule := func() {
		scheduleCancel()

		_, scheduleCancel = scheduleOnce(childCtx, op.Duration, func() {
			if try.Lock() {
				try.CancelAndUnlock()
				childCancel()
				op.Observable.Subscribe(ctx, sink)
			}
		})
	}

	doSchedule()

	source.Subscribe(childCtx, func(t Notification) {
		if try.Lock() {
			switch {
			case t.HasValue:
				sink(t)
				try.Unlock()
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
		op := timeoutOperator{timeout, Throw(ErrTimeout)}
		return source.Lift(op.Call)
	}
}
