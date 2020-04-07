package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/x/misc"
)

type windowWhenObservable struct {
	Source          rx.Observable
	ClosingSelector func() rx.Observable
}

func (obs windowWhenObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	type X struct {
		Window rx.Subject
	}
	window := rx.NewSubject()
	cx := make(chan *X, 1)
	cx <- &X{window}
	sink.Next(window.Observable)

	var (
		openWindow     func()
		avoidRecursive misc.AvoidRecursive
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
					t.Observe(x.Window.Observer)
					sink(t)
					return
				}
				x.Window.Complete()
				x.Window = rx.NewSubject()
				sink.Next(x.Window.Observable)
				cx <- x
				avoidRecursive.Do(openWindow)
			}
		}

		closingNotifier := obs.ClosingSelector()
		closingNotifier.Subscribe(ctx, observer.Notify)
	}

	avoidRecursive.Do(openWindow)

	if ctx.Err() != nil {
		return
	}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				t.Observe(x.Window.Observer)
				cx <- x
			default:
				close(cx)
				t.Observe(x.Window.Observer)
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
		obs := windowWhenObservable{source, closingSelector}
		return rx.Create(obs.Subscribe)
	}
}
