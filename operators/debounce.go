package operators

import (
	"context"
	"time"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/critical"
)

// Debounce emits a value from the source only after a particular time span,
// determined by another Observable, has passed without another source
// emission.
//
// It's like DebounceTime, but the time span of emission silence is determined
// by a second Observable.
func Debounce(durationSelector func(interface{}) rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return debounceObservable{source, durationSelector}.Subscribe
	}
}

// DebounceTime emits a value from the source only after a particular time
// span has passed without another source emission.
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

				ctx, cancel := context.WithCancel(ctx)
				scheduleCancel = cancel

				var observer rx.Observer

				observer = func(t rx.Notification) {
					observer = rx.Noop

					cancel()

					if critical.Enter(&x.Section) {
						switch {
						case t.HasValue:
							defer critical.Leave(&x.Section)

							if x.Latest.HasValue {
								value := x.Latest.Value

								x.Latest.Value = nil
								x.Latest.HasValue = false

								sink.Next(value)
							}

						case t.HasError:
							critical.Close(&x.Section)

							sink(t)

						default:
							critical.Leave(&x.Section)
						}
					}
				}

				obs1 := obs.DurationSelector(t.Value)

				obs1.Subscribe(ctx, observer.Sink)

			case t.HasError:
				critical.Close(&x.Section)

				sink(t)

			default:
				critical.Close(&x.Section)

				if x.Latest.HasValue {
					value := x.Latest.Value

					x.Latest.Value = nil
					x.Latest.HasValue = false

					sink.Next(value)
				}

				sink(t)
			}
		}
	})
}
