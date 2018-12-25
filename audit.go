package rx

import (
	"context"

	"github.com/b97tsk/rx/x/atomic"
)

type auditOperator struct {
	DurationSelector func(interface{}) Observable
}

func (op auditOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
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

	doSchedule := func(val interface{}) {
		if !state.Cas(stateHasValue, stateScheduled) {
			return
		}

		scheduleCtx, scheduleCancel := context.WithCancel(ctx)

		var observer Observer
		observer = func(t Notification) {
			observer = NopObserver
			scheduleCancel()
			if try.Lock() {
				if t.HasError {
					try.CancelAndUnlock()
					sink(t)
					return
				}
				sink.Next(latestValue)
				state.Store(stateZero)
				try.Unlock()
			}
		}

		obs := op.DurationSelector(val)
		obs.Subscribe(scheduleCtx, observer.Notify)
	}

	source.Subscribe(ctx, func(t Notification) {
		if try.Lock() {
			switch {
			case t.HasValue:
				latestValue = t.Value
				state.Cas(stateZero, stateHasValue)
				try.Unlock()
				doSchedule(t.Value)
			default:
				try.CancelAndUnlock()
				sink(t)
			}
		}
	})

	return ctx, cancel
}

// Audit ignores source values for a duration determined by another Observable,
// then emits the most recent value from the source Observable, then repeats
// this process.
//
// It's like AuditTime, but the silencing duration is determined by a second
// Observable.
func (Operators) Audit(durationSelector func(interface{}) Observable) OperatorFunc {
	return func(source Observable) Observable {
		op := auditOperator{durationSelector}
		return source.Lift(op.Call)
	}
}
