package operators

import (
	"context"

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
			observer = rx.NopObserver
			defer cancel()
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
		}

		obs := obs.DurationSelector(val)
		go obs.Subscribe(throttleCtx, observer.Notify)
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
