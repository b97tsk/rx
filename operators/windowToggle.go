package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/subject"
)

type windowToggleObservable struct {
	Source          rx.Observable
	Openings        rx.Observable
	ClosingSelector func(interface{}) rx.Observable
}

type windowToggleContext struct {
	Cancel context.CancelFunc
	Window *subject.Subject
}

func (obs windowToggleObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)
	sink = sink.WithCancel(cancel)

	type X struct {
		Contexts []*windowToggleContext
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	cleanupContexts := func(x *X, t rx.Notification) {
		for _, c := range x.Contexts {
			c.Cancel()
			c.Window.Sink(t)
		}
	}

	obs.Openings.Subscribe(ctx, func(t rx.Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				ctx, cancel := context.WithCancel(ctx)
				newContext := &windowToggleContext{
					Cancel: cancel,
					Window: subject.NewSubject(),
				}
				x.Contexts = append(x.Contexts, newContext)
				sink.Next(newContext.Window.Observable)

				cx <- x

				var observer rx.Observer
				observer = func(t rx.Notification) {
					observer = rx.Noop
					cancel()
					if x, ok := <-cx; ok {
						if t.HasError {
							close(cx)
							cleanupContexts(x, t)
							sink(t)
							return
						}
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
						cx <- x
					}
				}

				closingNotifier := obs.ClosingSelector(t.Value)
				closingNotifier.Subscribe(ctx, observer.Sink)

			case t.HasError:
				close(cx)
				cleanupContexts(x, t)
				sink(t)

			default:
				cx <- x
			}
		}
	})

	if ctx.Err() != nil {
		return
	}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				for _, c := range x.Contexts {
					c.Window.Sink(t)
				}

				cx <- x

			default:
				close(cx)
				cleanupContexts(x, t)
				sink(t)
			}
		}
	})
}

// WindowToggle branches out the source Observable values as a nested
// Observable starting from an emission from openings and ending when
// the output of closingSelector emits.
//
// It's like BufferToggle, but emits a nested Observable instead of a slice.
func WindowToggle(openings rx.Observable, closingSelector func(interface{}) rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return windowToggleObservable{source, openings, closingSelector}.Subscribe
	}
}
