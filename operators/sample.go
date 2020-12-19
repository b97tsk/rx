package operators

import (
	"context"
	"time"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/critical"
)

// Sample creates an Observable that emits the most recently emitted value from
// the source Observable whenever another Observable, the notifier, emits.
//
// It's like SampleTime, but samples whenever the notifier Observable emits
// something.
func Sample(notifier rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return sampleObservable{source, notifier}.Subscribe
	}
}

// SampleTime creates an Observable that emits the most recently emitted value
// from the source Observable within periodic time intervals.
func SampleTime(d time.Duration) rx.Operator {
	return Sample(rx.Ticker(d))
}

type sampleObservable struct {
	Source   rx.Observable
	Notifier rx.Observable
}

func (obs sampleObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel)

	var x struct {
		critical.Section
		Latest struct {
			Value    interface{}
			HasValue bool
		}
	}

	obs.Notifier.Subscribe(ctx, func(t rx.Notification) {
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
	})

	if ctx.Err() != nil {
		return
	}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		if critical.Enter(&x.Section) {
			switch {
			case t.HasValue:
				x.Latest.Value = t.Value
				x.Latest.HasValue = true

				critical.Leave(&x.Section)

			case t.HasError:
				fallthrough

			default:
				critical.Close(&x.Section)

				sink(t)
			}
		}
	})
}
