package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/critical"
)

// WindowToggle branches out the source values as a nested Observable
// starting from an emission from openings and ending when the output of
// closingSelector emits.
//
// It's like BufferToggle, but emits a nested Observable instead of a slice.
func WindowToggle(openings rx.Observable, closingSelector func(interface{}) rx.Observable) rx.Operator {
	return WindowToggleConfig{
		Openings:        openings,
		ClosingSelector: closingSelector,
	}.Make()
}

// A WindowToggleConfig is a configuration for WindowToggle.
type WindowToggleConfig struct {
	Openings        rx.Observable
	ClosingSelector func(interface{}) rx.Observable
	WindowFactory   rx.SubjectFactory
}

// Make creates an Operator from this configuration.
func (config WindowToggleConfig) Make() rx.Operator {
	if config.Openings == nil {
		panic("WindowToggle: Openings is nil")
	}

	if config.ClosingSelector == nil {
		panic("WindowToggle: ClosingSelector is nil")
	}

	if config.WindowFactory == nil {
		config.WindowFactory = rx.Multicast
	}

	return func(source rx.Observable) rx.Observable {
		return windowToggleObservable{source, config}.Subscribe
	}
}

type windowToggleObservable struct {
	Source rx.Observable
	WindowToggleConfig
}

type windowToggleContext struct {
	Cancel context.CancelFunc
	Window rx.Observer
}

func (obs windowToggleObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel)

	var x struct {
		critical.Section
		Contexts []*windowToggleContext
	}

	cleanupContexts := func(t rx.Notification) {
		for _, c := range x.Contexts {
			c.Cancel()
			c.Window.Sink(t)
		}
	}

	obs.Openings.Subscribe(ctx, func(t rx.Notification) {
		if critical.Enter(&x.Section) {
			switch {
			case t.HasValue:
				ctx, cancel := context.WithCancel(ctx)

				window := obs.WindowFactory()

				newContext := &windowToggleContext{
					Cancel: cancel,
					Window: window.Observer,
				}
				x.Contexts = append(x.Contexts, newContext)

				sink.Next(window.Observable)

				critical.Leave(&x.Section)

				var observer rx.Observer

				observer = func(t rx.Notification) {
					observer = rx.Noop

					cancel()

					if critical.Enter(&x.Section) {
						switch {
						case t.HasValue:
							defer critical.Leave(&x.Section)

							for i, c := range x.Contexts {
								if c == newContext {
									copy(x.Contexts[i:], x.Contexts[i+1:])

									n := len(x.Contexts)
									x.Contexts[n-1] = nil
									x.Contexts = x.Contexts[:n-1]

									newContext.Window.Complete()

									break
								}
							}

						case t.HasError:
							critical.Close(&x.Section)

							cleanupContexts(t)
							sink(t)

						default:
							critical.Leave(&x.Section)
						}
					}
				}

				closingNotifier := obs.ClosingSelector(t.Value)

				closingNotifier.Subscribe(ctx, observer.Sink)

			case t.HasError:
				critical.Close(&x.Section)

				cleanupContexts(t)
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
				for _, c := range x.Contexts {
					c.Window.Sink(t)
				}

				critical.Leave(&x.Section)

			case t.HasError:
				fallthrough

			default:
				critical.Close(&x.Section)

				cleanupContexts(t)
				sink(t)
			}
		}
	})
}
