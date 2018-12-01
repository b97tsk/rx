package rx

import (
	"context"
)

// ThrottleOperator is an operator type.
type ThrottleOperator struct {
	DurationSelector func(interface{}) Observable
	Leading          bool
	Trailing         bool
}

// MakeFunc creates an OperatorFunc from this operator.
func (op ThrottleOperator) MakeFunc() OperatorFunc {
	return MakeFunc(op.Call)
}

// Call invokes an execution of this operator.
func (op ThrottleOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	var (
		throttleCtx    = canceledCtx
		throttleCancel = nothingToDo

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
						defer try.Unlock()
						sink.Next(trailingValue)
						hasTrailingValue = false
						doThrottle(trailingValue)
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
				defer try.Unlock()

				trailingValue = t.Value
				hasTrailingValue = true

				select {
				case <-throttleCtx.Done():
				default:
					return
				}

				doThrottle(t.Value)

				if leading {
					sink(t)
					hasTrailingValue = false
				}

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
		op := ThrottleOperator{
			DurationSelector: durationSelector,
			Leading:          true,
			Trailing:         false,
		}
		return source.Lift(op.Call)
	}
}
