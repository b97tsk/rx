package rx

import (
	"context"
	"time"
)

// A ThrottleTimeConfigure is a configure for ThrottleTime.
type ThrottleTimeConfigure struct {
	Duration time.Duration
	Leading  bool
	Trailing bool
}

// MakeFunc creates an OperatorFunc from this type.
func (conf ThrottleTimeConfigure) MakeFunc() OperatorFunc {
	return MakeFunc(throttleTimeOperator(conf).Call)
}

type throttleTimeOperator ThrottleTimeConfigure

func (op throttleTimeOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	throttleCtx, _ := Done()

	sink = Finally(sink, cancel)

	type X struct {
		TrailingValue    interface{}
		HasTrailingValue bool
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	var doThrottle func()

	doThrottle = func() {
		throttleCtx, _ = scheduleOnce(ctx, op.Duration, func() {
			if op.Trailing {
				if x, ok := <-cx; ok {
					if x.HasTrailingValue {
						sink.Next(x.TrailingValue)
						x.HasTrailingValue = false
						doThrottle()
					}
					cx <- x
				}
			}
		})
	}

	source.Subscribe(ctx, func(t Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				x.TrailingValue = t.Value
				x.HasTrailingValue = true
				if isDone(throttleCtx) {
					doThrottle()
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

// ThrottleTime creates an Observable that emits a value from the source
// Observable, then ignores subsequent source values for a duration, then
// repeats this process.
//
// ThrottleTime lets a value pass, then ignores source values for the next
// duration time.
func (Operators) ThrottleTime(duration time.Duration) OperatorFunc {
	return func(source Observable) Observable {
		op := throttleTimeOperator{
			Duration: duration,
			Leading:  true,
			Trailing: false,
		}
		return source.Lift(op.Call)
	}
}
