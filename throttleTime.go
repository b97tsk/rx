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

// Use creates an Operator from this configure.
func (configure ThrottleTimeConfigure) Use() Operator {
	return func(source Observable) Observable {
		return throttleTimeObservable{source, configure}.Subscribe
	}
}

type throttleTimeObservable struct {
	Source Observable
	ThrottleTimeConfigure
}

func (obs throttleTimeObservable) Subscribe(parent context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx := NewContext(parent)

	sink = DoAtLast(sink, ctx.AtLast)

	type X struct {
		TrailingValue    interface{}
		HasTrailingValue bool
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	var (
		doThrottle  func()
		throttleCtx context.Context
	)

	doThrottle = func() {
		throttleCtx, _ = scheduleOnce(ctx, obs.Duration, func() {
			if obs.Trailing {
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

	obs.Source.Subscribe(ctx, func(t Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				x.TrailingValue = t.Value
				x.HasTrailingValue = true
				if throttleCtx == nil || throttleCtx.Err() != nil {
					doThrottle()
					if obs.Leading {
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

	return ctx, ctx.Cancel
}

// ThrottleTime creates an Observable that emits a value from the source
// Observable, then ignores subsequent source values for a duration, then
// repeats this process.
//
// ThrottleTime lets a value pass, then ignores source values for the next
// duration time.
func (Operators) ThrottleTime(duration time.Duration) Operator {
	return ThrottleTimeConfigure{
		Duration: duration,
		Leading:  true,
		Trailing: false,
	}.Use()
}
