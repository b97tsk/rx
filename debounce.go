package rx

import (
	"context"
)

type debounceOperator struct {
	source           Operator
	durationSelector func(interface{}) Observable
}

func (op debounceOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	scheduleCtx, scheduleCancel := canceledCtx, noopFunc
	try := cancellableLocker{}
	latestValue := interface{}(nil)

	doSchedule := func(val interface{}) {
		scheduleCancel()

		scheduleCtx, scheduleCancel = context.WithCancel(ctx)

		mutable := MutableObserver{}

		mutable.Observer = ObserverFunc(func(t Notification) {
			if try.Lock() {
				if t.HasError {
					try.CancelAndUnlock()
					ob.Error(t.Value.(error))
					cancel()
					return
				}
				defer try.Unlock()
				defer scheduleCancel()
				mutable.Observer = NopObserver
				ob.Next(latestValue)
			}
		})

		obsv := op.durationSelector(val)
		obsv.Subscribe(scheduleCtx, &mutable)
	}

	op.source.Call(ctx, ObserverFunc(func(t Notification) {
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
	}))

	return ctx, cancel
}

// Debounce creates an Observable that emits a value from the source Observable
// only after a particular time span determined by another Observable has
// passed without another source emission.
//
// It's like DebounceTime, but the time span of emission silence is determined
// by a second Observable.
func (o Observable) Debounce(durationSelector func(interface{}) Observable) Observable {
	op := debounceOperator{
		source:           o.Op,
		durationSelector: durationSelector,
	}
	return Observable{op}
}
