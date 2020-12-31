package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/critical"
)

// Window branches out the source values as a nested Observable whenever
// windowBoundaries emits.
//
// It's like Buffer, but emits a nested Observable instead of a slice.
func Window(windowBoundaries rx.Observable) rx.Operator {
	return WindowConfigure{WindowBoundaries: windowBoundaries}.Make()
}

// A WindowConfigure is a configure for Window.
type WindowConfigure struct {
	WindowBoundaries rx.Observable
	WindowFactory    rx.DoubleFactory
}

// Make creates an Operator from this configure.
func (configure WindowConfigure) Make() rx.Operator {
	if configure.WindowBoundaries == nil {
		panic("Window: WindowBoundaries is nil")
	}

	if configure.WindowFactory == nil {
		configure.WindowFactory = rx.Multicast
	}

	return func(source rx.Observable) rx.Observable {
		return windowObservable{source, configure}.Subscribe
	}
}

type windowObservable struct {
	Source rx.Observable
	WindowConfigure
}

func (obs windowObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel)

	var x struct {
		critical.Section
		Window rx.Observer
	}

	window := obs.WindowFactory()
	x.Window = window.Observer
	sink.Next(window.Observable)

	obs.WindowBoundaries.Subscribe(ctx, func(t rx.Notification) {
		if critical.Enter(&x.Section) {
			switch {
			case t.HasValue:
				defer critical.Leave(&x.Section)

				x.Window.Complete()

				window := obs.WindowFactory()
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
