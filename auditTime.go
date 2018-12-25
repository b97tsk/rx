package rx

import (
	"context"
	"time"

	"github.com/b97tsk/rx/x/atomic"
)

type auditTimeOperator struct {
	Duration time.Duration
}

func (op auditTimeOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	const (
		stateZero = iota
		stateHasValue
		stateScheduled
	)

	var (
		latestValue interface{}
		state       atomic.Uint32

		try cancellableLocker
	)

	doSchedule := func() {
		if !state.Cas(stateHasValue, stateScheduled) {
			return
		}
		scheduleOnce(ctx, op.Duration, func() {
			if try.Lock() {
				sink.Next(latestValue)
				state.Store(stateZero)
				try.Unlock()
			}
		})
	}

	source.Subscribe(ctx, func(t Notification) {
		if try.Lock() {
			switch {
			case t.HasValue:
				latestValue = t.Value
				state.Cas(stateZero, stateHasValue)
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

// AuditTime ignores source values for a duration, then emits the most recent
// value from the source Observable, then repeats this process.
//
// When it sees a source values, it ignores that plus the next ones for a
// duration, and then it emits the most recent value from the source.
func (Operators) AuditTime(duration time.Duration) OperatorFunc {
	return func(source Observable) Observable {
		op := auditTimeOperator{duration}
		return source.Lift(op.Call)
	}
}
