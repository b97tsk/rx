package rx

import (
	"context"
)

// A ThrottleConfigure is a configure for Throttle.
type ThrottleConfigure struct {
	DurationSelector func(interface{}) Observable
	Leading          bool
	Trailing         bool
}

// MakeFunc creates an OperatorFunc from this type.
func (conf ThrottleConfigure) MakeFunc() OperatorFunc {
	return MakeFunc(throttleOperator(conf).Call)
}

type throttleOperator ThrottleConfigure

func (op throttleOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	throttleCtx, throttleCancel := Done()

	sink = Finally(sink, cancel)

	var (
		trailingValue    interface{}
		hasTrailingValue bool

		try cancellableLocker
	)

	leading, trailing := op.Leading, op.Trailing
	if !leading && !trailing {
		leading = true
	}

	var doThrottle func(interface{})

	doThrottle = func(val interface{}) {
		throttleCtx, throttleCancel = context.WithCancel(ctx)

		var observer Observer
		observer = func(t Notification) {
			observer = NopObserver
			defer throttleCancel()
			if trailing || t.HasError {
				if try.Lock() {
					switch {
					case t.HasError:
						try.CancelAndUnlock()
						sink(t)
					case hasTrailingValue:
						sink.Next(trailingValue)
						hasTrailingValue = false
						doThrottle(trailingValue)
						try.Unlock()
					}
				}
			}
		}

		obs := op.DurationSelector(val)
		go obs.Subscribe(throttleCtx, observer.Notify)
	}

	source.Subscribe(ctx, func(t Notification) {
		if try.Lock() {
			switch {
			case t.HasValue:
				trailingValue = t.Value
				hasTrailingValue = true
				if isDone(throttleCtx) {
					doThrottle(t.Value)
					if leading {
						sink(t)
						hasTrailingValue = false
					}
				}
				try.Unlock()

			default:
				try.CancelAndUnlock()
				if hasTrailingValue {
					sink.Next(trailingValue)
				}
				sink(t)
			}
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
func (Operators) Throttle(durationSelector func(interface{}) Observable) OperatorFunc {
	return func(source Observable) Observable {
		op := throttleOperator{
			DurationSelector: durationSelector,
			Leading:          true,
			Trailing:         false,
		}
		return source.Lift(op.Call)
	}
}
