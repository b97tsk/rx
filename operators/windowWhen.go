package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/critical"
	"github.com/b97tsk/rx/internal/norec"
)

// WindowWhen branches out the source Observable values as a nested Observable
// using a factory function of closing Observables to determine when to start
// a new window.
//
// It's like BufferWhen, but emits a nested Observable instead of a slice.
func WindowWhen(closingSelector func() rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return windowWhenObservable{source, closingSelector}.Subscribe
	}
}

type windowWhenObservable struct {
	Source          rx.Observable
	ClosingSelector func() rx.Observable
}

func (obs windowWhenObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel)

	var x struct {
		critical.Section
		Window rx.Observer
	}

	window := rx.Multicast()
	x.Window = window.Observer
	sink.Next(window.Observable)

	var openWindow func()

	openWindow = norec.Wrap(func() {
		if ctx.Err() != nil {
			return
		}

		ctx, cancel := context.WithCancel(ctx)

		var observer rx.Observer

		observer = func(t rx.Notification) {
			observer = rx.Noop

			cancel()

			if critical.Enter(&x.Section) {
				switch {
				case t.HasValue:
					x.Window.Complete()

					window := rx.Multicast()
					x.Window = window.Observer
					sink.Next(window.Observable)

					critical.Leave(&x.Section)

					openWindow()

				case t.HasError:
					critical.Close(&x.Section)

					x.Window.Sink(t)
					sink(t)

				default:
					critical.Leave(&x.Section)
				}
			}
		}

		closingNotifier := obs.ClosingSelector()

		closingNotifier.Subscribe(ctx, observer.Sink)
	})

	openWindow()

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
