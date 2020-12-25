package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/critical"
)

// Window branches out the source Observable values as a nested Observable
// whenever windowBoundaries emits.
//
// It's like Buffer, but emits a nested Observable instead of a slice.
func Window(windowBoundaries rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return windowObservable{source, windowBoundaries}.Subscribe
	}
}

type windowObservable struct {
	Source           rx.Observable
	WindowBoundaries rx.Observable
}

func (obs windowObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel)

	var x struct {
		critical.Section
		Window rx.Observer
	}

	window := rx.Multicast()
	x.Window = window.Observer
	sink.Next(window.Observable)

	obs.WindowBoundaries.Subscribe(ctx, func(t rx.Notification) {
		if critical.Enter(&x.Section) {
			switch {
			case t.HasValue:
				defer critical.Leave(&x.Section)

				x.Window.Complete()

				window := rx.Multicast()
				x.Window = window.Observer
				sink.Next(window.Observable)

			case t.HasError:
				fallthrough

			default:
				critical.Close(&x.Section)

				x.Window.Sink(t)
				sink(t)
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
				x.Window.Sink(t)

				critical.Leave(&x.Section)

			case t.HasError:
				fallthrough

			default:
				critical.Close(&x.Section)

				x.Window.Sink(t)
				sink(t)
			}
		}
	})
}
