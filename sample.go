package rx

import (
	"context"
	"time"

	"github.com/b97tsk/rx/internal/critical"
)

// SampleTime emits the most recently emitted value from the source Observalbe
// within periodic time intervals.
func SampleTime[T any](d time.Duration) Operator[T, T] {
	return Sample[T](Ticker(d))
}

// Sample emits the most recently emitted value from the source Observable
// whenever notifier, another Observable, emits a value.
func Sample[T, U any](notifier Observable[U]) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return sampleObservable[T, U]{source, notifier}.Subscribe
		},
	)
}

type sampleObservable[T, U any] struct {
	Source   Observable[T]
	Notifier Observable[U]
}

func (obs sampleObservable[T, U]) Subscribe(ctx context.Context, sink Observer[T]) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.OnLastNotification(cancel)

	var x struct {
		critical.Section
		Latest struct {
			Value    T
			HasValue bool
		}
	}

	obs.Notifier.Subscribe(ctx, func(n Notification[U]) {
		if critical.Enter(&x.Section) {
			switch {
			case n.HasValue:
				if x.Latest.HasValue {
					sink.Next(x.Latest.Value)
					x.Latest.HasValue = false
				}

				critical.Leave(&x.Section)

			case n.HasError:
				critical.Close(&x.Section)

				sink.Error(n.Error)

			default:
				critical.Leave(&x.Section)
			}
		}
	})

	if err := getErr(ctx); err != nil {
		if critical.Enter(&x.Section) {
			critical.Close(&x.Section)
			sink.Error(err)
		}

		return
	}

	obs.Source.Subscribe(ctx, func(n Notification[T]) {
		if critical.Enter(&x.Section) {
			switch {
			case n.HasValue:
				x.Latest.Value = n.Value
				x.Latest.HasValue = true

				critical.Leave(&x.Section)

			default:
				critical.Close(&x.Section)

				sink(n)
			}
		}
	})
}
