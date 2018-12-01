package rx

import (
	"context"
	"time"
)

// ThrottleTimeOperator is an operator type.
type ThrottleTimeOperator struct {
	Duration time.Duration
	Leading  bool
	Trailing bool
}

// MakeFunc creates an OperatorFunc from this operator.
func (op ThrottleTimeOperator) MakeFunc() OperatorFunc {
	return MakeFunc(op.Call)
}

// Call invokes an execution of this operator.
func (op ThrottleTimeOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	var (
		throttleCtx  = canceledCtx

		trailingValue    interface{}
		hasTrailingValue bool

		try cancellableLocker
	)

	leading, trailing := op.Leading, op.Trailing
	if !leading && !trailing {
		leading = true
	}

	var doThrottle func()

	doThrottle = func() {
		throttleCtx, _ = scheduleOnce(ctx, op.Duration, func() {
			if trailing {
				if try.Lock() {
					defer try.Unlock()
					if hasTrailingValue {
						sink.Next(trailingValue)
						hasTrailingValue = false
						doThrottle()
					}
				}
			}
		})
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

				doThrottle()

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

// ThrottleTime creates an Observable that emits a value from the source
// Observable, then ignores subsequent source values for a duration, then
// repeats this process.
//
// ThrottleTime lets a value pass, then ignores source values for the next
// duration time.
func (Operators) ThrottleTime(duration time.Duration) OperatorFunc {
	return func(source Observable) Observable {
		op := ThrottleTimeOperator{
			Duration: duration,
			Leading:  true,
			Trailing: false,
		}
		return source.Lift(op.Call)
	}
}
