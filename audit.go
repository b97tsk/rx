package rx

import (
	"context"
)

type auditOperator struct {
	DurationSelector func(interface{}) Observable
}

func (op auditOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	scheduleCtx, scheduleCancel := canceledCtx, doNothing
	scheduleDone := scheduleCtx.Done()

	var (
		latestValue interface{}
		try         cancellableLocker
	)

	doSchedule := func(val interface{}) {
		select {
		case <-scheduleDone:
		default:
			return
		}

		scheduleCtx, scheduleCancel = context.WithCancel(ctx)
		scheduleDone = scheduleCtx.Done()

		var observer Observer
		observer = func(t Notification) {
			observer = NopObserver
			scheduleCancel()
			if try.Lock() {
				if t.HasError {
					try.CancelAndUnlock()
					sink(t)
					cancel()
					return
				}
				defer try.Unlock()
				sink.Next(latestValue)
			}
		}

		obsv := op.DurationSelector(val)
		obsv.Subscribe(scheduleCtx, observer.Notify)
	}

	source.Subscribe(ctx, func(t Notification) {
		if try.Lock() {
			switch {
			case t.HasValue:
				latestValue = t.Value
				try.Unlock()
				doSchedule(t.Value)
			default:
				try.CancelAndUnlock()
				sink(t)
				cancel()
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
func (o Observable) Audit(durationSelector func(interface{}) Observable) Observable {
	op := auditOperator{durationSelector}
	return o.Lift(op.Call)
}
