package rx

import (
	"context"
)

type debounceOperator struct {
	durationSelector func(interface{}) Observable
}

func (op debounceOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	scheduleCtx, scheduleCancel := canceledCtx, doNothing

	var (
		latestValue interface{}
		try         cancellableLocker
	)

	doSchedule := func(val interface{}) {
		scheduleCancel()

		scheduleCtx, scheduleCancel = context.WithCancel(ctx)

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

		obsv := op.durationSelector(val)
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

// Debounce creates an Observable that emits a value from the source Observable
// only after a particular time span, determined by another Observable, has
// passed without another source emission.
//
// It's like DebounceTime, but the time span of emission silence is determined
// by a second Observable.
func (o Observable) Debounce(durationSelector func(interface{}) Observable) Observable {
	op := debounceOperator{durationSelector}
	return o.Lift(op.Call)
}
