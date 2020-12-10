package operators

import (
	"context"
	"time"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/critical"
)

// Debounce creates an Observable that emits a value from the source Observable
// only after a particular time span, determined by another Observable, has
// passed without another source emission.
//
// It's like DebounceTime, but the time span of emission silence is determined
// by a second Observable.
func Debounce(durationSelector func(interface{}) rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return debounceObservable{source, durationSelector}.Subscribe
	}
}

// DebounceTime creates an Observable that emits a value from the source
// Observable only after a particular time span has passed without another
// source emission.
func DebounceTime(d time.Duration) rx.Operator {
	obsTimer := rx.Timer(d)

	durationSelector := func(interface{}) rx.Observable { return obsTimer }

	return Debounce(durationSelector)
}

type debounceObservable struct {
	Source           rx.Observable
	DurationSelector func(interface{}) rx.Observable
}

func (obs debounceObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel)

	var x struct {
		critical.Section
		Latest struct {
			Value    interface{}
			HasValue bool
		}
	}

	var scheduleCancel context.CancelFunc

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		if critical.Enter(&x.Section) {
			switch {
			case t.HasValue:
				x.Latest.Value = t.Value
				x.Latest.HasValue = true

				critical.Leave(&x.Section)

				if scheduleCancel != nil {
					scheduleCancel()
				}

				var scheduleCtx context.Context

				scheduleCtx, scheduleCancel = context.WithCancel(ctx)

				var observer rx.Observer

				observer = func(t rx.Notification) {
					observer = rx.Noop

					scheduleCancel()

					if critical.Enter(&x.Section) {
						if t.HasError {
							critical.Close(&x.Section)
							sink(t)
							return
						}

						if x.Latest.HasValue {
							x.Latest.HasValue = false
							sink.Next(x.Latest.Value)
						}

						critical.Leave(&x.Section)
					}
				}

				obs1 := obs.DurationSelector(t.Value)
				obs1.Subscribe(scheduleCtx, observer.Sink)

			case t.HasError:
				critical.Close(&x.Section)

				sink(t)

			default:
				critical.Close(&x.Section)

				if x.Latest.HasValue {
					sink.Next(x.Latest.Value)
				}

				sink(t)
			}
		}
	})
}
