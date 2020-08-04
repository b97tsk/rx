package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/misc"
)

type windowWhenObservable struct {
	Source          rx.Observable
	ClosingSelector func() rx.Observable
}

func (obs windowWhenObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)
	sink = sink.WithCancel(cancel)

	type X struct {
		Window rx.Observer
	}
	window := rx.Multicast()
	cx := make(chan *X, 1)
	cx <- &X{window.Observer}
	sink.Next(window.Observable)

	var (
		openWindow     func()
		avoidRecursion misc.AvoidRecursion
	)

	openWindow = func() {
		if ctx.Err() != nil {
			return
		}

		ctx, cancel := context.WithCancel(ctx)

		var observer rx.Observer
		observer = func(t rx.Notification) {
			observer = rx.Noop
			cancel()
			if x, ok := <-cx; ok {
				if t.HasError {
					close(cx)
					x.Window.Sink(t)
					sink(t)
					return
				}
				x.Window.Complete()
				window := rx.Multicast()
				x.Window = window.Observer
				sink.Next(window.Observable)
				cx <- x
				avoidRecursion.Do(openWindow)
			}
		}

		closingNotifier := obs.ClosingSelector()
		closingNotifier.Subscribe(ctx, observer.Sink)
	}

	avoidRecursion.Do(openWindow)

	if ctx.Err() != nil {
		return
	}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				x.Window.Sink(t)
				cx <- x
			default:
				close(cx)
				x.Window.Sink(t)
				sink(t)
			}
		}
	})
}

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
