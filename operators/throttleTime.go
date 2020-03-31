package operators

import (
	"context"
	"time"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/x/schedule"
)

// A ThrottleTimeConfigure is a configure for ThrottleTime.
type ThrottleTimeConfigure struct {
	Duration time.Duration
	Leading  bool
	Trailing bool
}

// Use creates an Operator from this configure.
func (configure ThrottleTimeConfigure) Use() rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := throttleTimeObservable{source, configure}
		return rx.Create(obs.Subscribe)
	}
}

type throttleTimeObservable struct {
	Source rx.Observable
	ThrottleTimeConfigure
}

func (obs throttleTimeObservable) Subscribe(ctx context.Context, sink rx.Observer) {
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
		throttleCtx, _ = schedule.ScheduleOnce(ctx, obs.Duration, func() {
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

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
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
}

// ThrottleTime creates an Observable that emits a value from the source
// Observable, then ignores subsequent source values for a duration, then
// repeats this process.
//
// ThrottleTime lets a value pass, then ignores source values for the next
// duration time.
func ThrottleTime(duration time.Duration) rx.Operator {
	return ThrottleTimeConfigure{
		Duration: duration,
		Leading:  true,
		Trailing: false,
	}.Use()
}
