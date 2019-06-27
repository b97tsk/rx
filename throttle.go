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

	type X struct {
		TrailingValue    interface{}
		HasTrailingValue bool
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	var doThrottle func(interface{})

	doThrottle = func(val interface{}) {
		throttleCtx, throttleCancel = context.WithCancel(ctx)

		var observer Observer
		observer = func(t Notification) {
			observer = NopObserver
			defer throttleCancel()
			if op.Trailing || t.HasError {
				if x, ok := <-cx; ok {
					switch {
					case t.HasError:
						close(cx)
						sink(t)
					case x.HasTrailingValue:
						sink.Next(x.TrailingValue)
						x.HasTrailingValue = false
						doThrottle(x.TrailingValue)
						cx <- x
					}
				}
			}
		}

		obs := op.DurationSelector(val)
		go obs.Subscribe(throttleCtx, observer.Notify)
	}

	source.Subscribe(ctx, func(t Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				x.TrailingValue = t.Value
				x.HasTrailingValue = true
				if isDone(throttleCtx) {
					doThrottle(t.Value)
					if op.Leading {
						sink(t)
						x.HasTrailingValue = false
					}
				}
				cx <- x

			default:
				close(cx)
				if x.HasTrailingValue {
					sink.Next(x.TrailingValue)
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
