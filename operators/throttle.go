package operators

import (
	"context"
	"time"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/critical"
)

// Throttle emits a value from the source, then ignores subsequent source
// values for a duration determined by another Observable, then repeats this
// process.
//
// It's like ThrottleTime, but the silencing duration is determined by a second
// Observable.
func Throttle(durationSelector func(interface{}) rx.Observable) rx.Operator {
	return ThrottleConfig{
		DurationSelector: durationSelector,
		Leading:          true,
		Trailing:         false,
	}.Make()
}

// ThrottleTime emits a value from the source, then ignores subsequent source
// values for a duration, then repeats this process.
//
// ThrottleTime lets a value pass, then ignores source values for the next
// duration time.
func ThrottleTime(d time.Duration) rx.Operator {
	return ThrottleTimeConfig{
		Duration: d,
		Leading:  true,
		Trailing: false,
	}.Make()
}

// A ThrottleConfig is a configuration for Throttle.
type ThrottleConfig struct {
	DurationSelector func(interface{}) rx.Observable
	Leading          bool
	Trailing         bool
}

// Make creates an Operator from this configuration.
func (config ThrottleConfig) Make() rx.Operator {
	if config.DurationSelector == nil {
		panic("Throttle: DurationSelector is nil")
	}

	return func(source rx.Observable) rx.Observable {
		return throttleObservable{source, config}.Subscribe
	}
}

// A ThrottleTimeConfig is a configuration for ThrottleTime.
type ThrottleTimeConfig struct {
	Duration time.Duration
	Leading  bool
	Trailing bool
}

// Make creates an Operator from this configuration.
func (config ThrottleTimeConfig) Make() rx.Operator {
	obsTimer := rx.Timer(config.Duration)

	durationSelector := func(interface{}) rx.Observable { return obsTimer }

	return ThrottleConfig{
		DurationSelector: durationSelector,
		Leading:          config.Leading,
		Trailing:         config.Trailing,
	}.Make()
}

type throttleObservable struct {
	Source rx.Observable
	ThrottleConfig
}

func (obs throttleObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel)

	var x struct {
		critical.Section
		Trailing struct {
			Value    interface{}
			HasValue bool
		}
		Throttling bool
		Completed  bool
	}

	var doThrottle func(interface{})

	doThrottle = func(val interface{}) {
		x.Throttling = true

		throttleCtx, throttleCancel := context.WithCancel(ctx)

		var observer rx.Observer

		observer = func(t rx.Notification) {
			observer = rx.Noop

			throttleCancel()

			if critical.Enter(&x.Section) {
				x.Throttling = false

				switch {
				case t.HasValue:
					defer critical.Leave(&x.Section)

					if obs.Trailing && x.Trailing.HasValue {
						value := x.Trailing.Value

						x.Trailing.Value = nil
						x.Trailing.HasValue = false

						sink.Next(value)

						if !x.Completed {
							doThrottle(value)
						}
					}

					if x.Completed {
						sink.Complete()
					}

				case t.HasError:
					critical.Close(&x.Section)

					sink(t)

				default:
					defer critical.Leave(&x.Section)

					if x.Completed {
						sink(t)
					}
				}
			}
		}

		obs1 := obs.DurationSelector(val)

		go obs1.Subscribe(throttleCtx, observer.Sink)
	}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		if critical.Enter(&x.Section) {
			switch {
			case t.HasValue:
				defer critical.Leave(&x.Section)

				x.Trailing.Value = t.Value
				x.Trailing.HasValue = true

				if !x.Throttling {
					doThrottle(t.Value)

					if obs.Leading {
						x.Trailing.Value = nil
						x.Trailing.HasValue = false

						sink(t)
					}
				}

			case t.HasError:
				critical.Close(&x.Section)

				sink(t)

			default:
				x.Completed = true

				if obs.Trailing && x.Trailing.HasValue && x.Throttling {
					critical.Leave(&x.Section)
					break
				}

				critical.Close(&x.Section)

				sink(t)
			}
		}
	})
}
