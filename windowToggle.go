package rx

import (
	"context"
)

type windowToggleOperator struct {
	Openings        Observable
	ClosingSelector func(interface{}) Observable
}

type windowToggleContext struct {
	Cancel context.CancelFunc
	Window Subject
}

func (op windowToggleOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	type X struct {
		Contexts []*windowToggleContext
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	cleanupContexts := func(x *X, t Notification) {
		for _, c := range x.Contexts {
			c.Cancel()
			t.Observe(c.Window.Observer)
		}
	}

	op.Openings.Subscribe(ctx, func(t Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				ctx, cancel := context.WithCancel(ctx)
				newContext := &windowToggleContext{
					Cancel: cancel,
					Window: NewSubject(),
				}
				x.Contexts = append(x.Contexts, newContext)
				sink.Next(newContext.Window.Observable)

				cx <- x

				var observer Observer
				observer = func(t Notification) {
					observer = NopObserver
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

				closingNotifier := op.ClosingSelector(t.Value)
				closingNotifier.Subscribe(ctx, observer.Notify)

			case t.HasError:
				close(cx)
				cleanupContexts(x, t)
				sink(t)

			default:
				cx <- x
			}
		}
	})

	if isDone(ctx) {
		return Done()
	}

	source.Subscribe(ctx, func(t Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				for _, c := range x.Contexts {
					t.Observe(c.Window.Observer)
				}

				cx <- x

			default:
				close(cx)
				cleanupContexts(x, t)
				sink(t)
			}
		}
	})

	return ctx, cancel
}

// WindowToggle branches out the source Observable values as a nested
// Observable starting from an emission from openings and ending when
// the output of closingSelector emits.
//
// It's like BufferToggle, but emits a nested Observable instead of a slice.
func (Operators) WindowToggle(openings Observable, closingSelector func(interface{}) Observable) OperatorFunc {
	return func(source Observable) Observable {
		op := windowToggleOperator{openings, closingSelector}
		return source.Lift(op.Call)
	}
}
