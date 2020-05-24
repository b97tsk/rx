package operators

import (
	"context"
	"time"

	"github.com/b97tsk/rx"
)

// A ThrottleConfigure is a configure for Throttle.
type ThrottleConfigure struct {
	DurationSelector func(interface{}) rx.Observable
	Leading          bool
	Trailing         bool
}

// Use creates an Operator from this configure.
func (configure ThrottleConfigure) Use() rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := throttleObservable{source, configure}
		return rx.Create(obs.Subscribe)
	}
}

// A ThrottleTimeConfigure is a configure for ThrottleTime.
type ThrottleTimeConfigure struct {
	Duration time.Duration
	Leading  bool
	Trailing bool
}

// Use creates an Operator from this configure.
func (configure ThrottleTimeConfigure) Use() rx.Operator {
	durationSelector := func(interface{}) rx.Observable {
		return rx.Timer(configure.Duration)
	}
	return ThrottleConfigure{
		DurationSelector: durationSelector,
		Leading:          configure.Leading,
		Trailing:         configure.Trailing,
	}.Use()
}

type throttleObservable struct {
	Source rx.Observable
	ThrottleConfigure
}

func (obs throttleObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	type X struct {
		TrailingValue    interface{}
		HasTrailingValue bool
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	var (
		doThrottle  func(interface{})
		throttleCtx context.Context
	)

	doThrottle = func(val interface{}) {
		ctx, cancel := context.WithCancel(ctx)
		throttleCtx = ctx

		var observer rx.Observer
		observer = func(t rx.Notification) {
			observer = rx.Noop
			if obs.Trailing || t.HasError {
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
			cancel()
		}

		obs := obs.DurationSelector(val)
		go obs.Subscribe(throttleCtx, observer.Sink)
	}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				x.TrailingValue = t.Value
				x.HasTrailingValue = true
				if throttleCtx == nil || throttleCtx.Err() != nil {
					doThrottle(t.Value)
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

// Throttle creates an Observable that emits a value from the source
// Observable, then ignores subsequent source values for a duration determined
// by another Observable, then repeats this process.
//
// It's like ThrottleTime, but the silencing duration is determined by a second
// Observable.
func Throttle(durationSelector func(interface{}) rx.Observable) rx.Operator {
	return ThrottleConfigure{
		DurationSelector: durationSelector,
		Leading:          true,
		Trailing:         false,
	}.Use()
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
