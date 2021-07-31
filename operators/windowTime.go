package operators

import (
	"context"
	"time"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/critical"
)

// WindowTime branches out the source values as a nested Observable
// periodically in time.
//
// It's like BufferTime, but emits a nested Observable instead of a slice.
func WindowTime(d time.Duration) rx.Operator {
	return WindowTimeConfig{TimeSpan: d}.Make()
}

// A WindowTimeConfig is a configuration for WindowTime.
type WindowTimeConfig struct {
	TimeSpan         time.Duration
	CreationInterval time.Duration
	MaxWindowSize    int
	WindowFactory    rx.SubjectFactory
}

// Make creates an Operator from this configuration.
func (config WindowTimeConfig) Make() rx.Operator {
	if config.WindowFactory == nil {
		config.WindowFactory = rx.Multicast
	}

	return func(source rx.Observable) rx.Observable {
		return windowTimeObservable{source, config}.Subscribe
	}
}

type windowTimeObservable struct {
	Source rx.Observable
	WindowTimeConfig
}

type windowTimeContext struct {
	Cancel context.CancelFunc
	Window rx.Observer
	Size   int
}

func (obs windowTimeObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel)

	var x struct {
		critical.Section
		Contexts []*windowTimeContext
	}

	var closeContext func(*windowTimeContext)

	obsTimer := rx.Timer(obs.TimeSpan)

	openContextLocked := func() {
		ctx, cancel := context.WithCancel(ctx)

		window := obs.WindowFactory()

		newContext := &windowTimeContext{
			Cancel: cancel,
			Window: window.Observer,
		}
		x.Contexts = append(x.Contexts, newContext)

		obsTimer.Subscribe(ctx, func(t rx.Notification) {
			if t.HasValue {
				return
			}

			closeContext(newContext)
		})

		sink.Next(window.Observable)
	}

	openContext := func() {
		if critical.Enter(&x.Section) {
			defer critical.Leave(&x.Section)

			openContextLocked()
		}
	}

	closeContext = func(toBeClosed *windowTimeContext) {
		toBeClosed.Cancel()

		if critical.Enter(&x.Section) {
			defer critical.Leave(&x.Section)

			for i, c := range x.Contexts {
				if c == toBeClosed {
					copy(x.Contexts[i:], x.Contexts[i+1:])

					n := len(x.Contexts)
					x.Contexts[n-1] = nil
					x.Contexts = x.Contexts[:n-1]

					toBeClosed.Window.Complete()

					if obs.CreationInterval <= 0 {
						openContextLocked()
					}

					break
				}
			}
		}
	}

	openContext()

	if obs.CreationInterval > 0 {
		rx.Ticker(obs.CreationInterval).Subscribe(ctx, func(t rx.Notification) {
			if t.HasValue {
				openContext()
			}
		})
	}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		if critical.Enter(&x.Section) {
			switch {
			case t.HasValue:
				var windowFullContexts []*windowTimeContext

				for _, c := range x.Contexts {
					c.Size++

					c.Window.Sink(t)

					if c.Size == obs.MaxWindowSize {
						windowFullContexts = append(windowFullContexts, c)
					}
				}

				critical.Leave(&x.Section)

				for _, c := range windowFullContexts {
					closeContext(c)
				}

			case t.HasError:
				fallthrough

			default:
				critical.Close(&x.Section)

				for _, c := range x.Contexts {
					c.Cancel()
					c.Window.Sink(t)
				}

				sink(t)
			}
		}
	})
}
