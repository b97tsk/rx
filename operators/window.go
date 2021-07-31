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
	return WindowConfig{WindowBoundaries: windowBoundaries}.Make()
}

// A WindowConfig is a configuration for Window.
type WindowConfig struct {
	WindowBoundaries rx.Observable
	WindowFactory    rx.SubjectFactory
}

// Make creates an Operator from this configuration.
func (config WindowConfig) Make() rx.Operator {
	if config.WindowBoundaries == nil {
		panic("Window: WindowBoundaries is nil")
	}

	if config.WindowFactory == nil {
		config.WindowFactory = rx.Multicast
	}

	return func(source rx.Observable) rx.Observable {
		return windowObservable{source, config}.Subscribe
	}
}

type windowObservable struct {
	Source rx.Observable
	WindowConfig
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
