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

	_, scheduleCancel := Done()
	childCtx, childCancel := context.WithCancel(ctx)

	type X struct{}
	cx := make(chan X, 1)
	cx <- X{}

	doSchedule := func() {
		scheduleCancel()

		_, scheduleCancel = scheduleOnce(childCtx, op.Duration, func() {
			if _, ok := <-cx; ok {
				close(cx)
				childCancel()
				op.Observable.Subscribe(ctx, sink)
			}
		})
	}

	doSchedule()

	source.Subscribe(childCtx, func(t Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				sink(t)
				cx <- x
				doSchedule()
			default:
				close(cx)
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
