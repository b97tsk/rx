package rx

import (
	"context"
)

type debounceOperator struct {
	DurationSelector func(interface{}) Observable
}

func (op debounceOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	var (
		scheduleCtx    = canceledCtx
		scheduleCancel = nothingToDo

		latestValue    interface{}
		hasLatestValue bool

		try cancellableLocker
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
					return
				}
				if hasLatestValue {
					sink.Next(latestValue)
					hasLatestValue = false
				}
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
				hasLatestValue = true
				try.Unlock()
				doSchedule(t.Value)

			case t.HasError:
				try.CancelAndUnlock()
				sink(t)

			default:
				try.CancelAndUnlock()
				if hasLatestValue {
					sink.Next(latestValue)
				}
				sink(t)
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
func (Operators) Debounce(durationSelector func(interface{}) Observable) OperatorFunc {
	return func(source Observable) Observable {
		op := debounceOperator{durationSelector}
		return source.Lift(op.Call)
	}
}
