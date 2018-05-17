package rx

import (
	"context"
)

type auditOperator struct {
	durationSelector func(interface{}) Observable
}

func (op auditOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	scheduleCtx, scheduleCancel := canceledCtx, noopFunc
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

		var mutableObserver Observer

		mutableObserver = func(t Notification) {
			if try.Lock() {
				if t.HasError {
					try.CancelAndUnlock()
					ob.Error(t.Value.(error))
					cancel()
					return
				}
				defer try.Unlock()
				defer scheduleCancel()
				mutableObserver = NopObserver
				ob.Next(latestValue)
			}
		}

		obsv := op.durationSelector(val)
		obsv.Subscribe(scheduleCtx, func(t Notification) { t.Observe(mutableObserver) })
	}

	source.Subscribe(ctx, func(t Notification) {
		if try.Lock() {
			switch {
			case t.HasValue:
				latestValue = t.Value
				try.Unlock()
				doSchedule(t.Value)
			case t.HasError:
				try.CancelAndUnlock()
				ob.Error(t.Value.(error))
				cancel()
			default:
				try.CancelAndUnlock()
				ob.Complete()
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
