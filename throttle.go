package rx

import (
	"context"
)

type throttleOperator struct {
	durationSelector func(interface{}) Observable
}

func (op throttleOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	scheduleCtx, scheduleCancel := canceledCtx, noopFunc
	scheduleDone := scheduleCtx.Done()

	source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			select {
			case <-scheduleDone:
			default:
				return
			}

			ob.Next(t.Value)

			scheduleCtx, scheduleCancel = context.WithCancel(ctx)
			scheduleDone = scheduleCtx.Done()

			var mutableObserver Observer

			mutableObserver = func(t Notification) {
				if t.HasError {
					ob.Error(t.Value.(error))
					cancel()
					return
				}
				mutableObserver = NopObserver
				scheduleCancel()
			}

			obsv := op.durationSelector(t.Value)
			obsv.Subscribe(scheduleCtx, func(t Notification) { t.Observe(mutableObserver) })

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

// Throttle creates an Observable that emits a value from the source
// Observable, then ignores subsequent source values for a duration determined
// by another Observable, then repeats this process.
//
// It's like ThrottleTime, but the silencing duration is determined by a second
// Observable.
func (o Observable) Throttle(durationSelector func(interface{}) Observable) Observable {
	op := throttleOperator{durationSelector}
	return o.Lift(op.Call)
}
