package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/critical"
)

// BufferToggle buffers the source Observable values starting from an emission
// from openings and ending when the output of closingSelector emits.
//
// BufferToggle collects values from the past as a slice, starts collecting
// only when opening emits, and calls the closingSelector function to get an
// Observable that tells when to close the buffer.
func BufferToggle(openings rx.Observable, closingSelector func(interface{}) rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return bufferToggleObservable{source, openings, closingSelector}.Subscribe
	}
}

type bufferToggleObservable struct {
	Source          rx.Observable
	Openings        rx.Observable
	ClosingSelector func(interface{}) rx.Observable
}

type bufferToggleContext struct {
	Cancel context.CancelFunc
	Buffer []interface{}
}

func (obs bufferToggleObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel)

	var x struct {
		critical.Section
		Contexts []*bufferToggleContext
	}

	obs.Openings.Subscribe(ctx, func(t rx.Notification) {
		if critical.Enter(&x.Section) {
			switch {
			case t.HasValue:
				ctx, cancel := context.WithCancel(ctx)

				newContext := &bufferToggleContext{Cancel: cancel}
				x.Contexts = append(x.Contexts, newContext)

				critical.Leave(&x.Section)

				var observer rx.Observer
				observer = func(t rx.Notification) {
					observer = rx.Noop
					cancel()

					if critical.Enter(&x.Section) {
						if t.HasError {
							critical.Close(&x.Section)
							sink(t)
							return
						}

						for i, c := range x.Contexts {
							if c == newContext {
								copy(x.Contexts[i:], x.Contexts[i+1:])

								n := len(x.Contexts)
								x.Contexts[n-1] = nil
								x.Contexts = x.Contexts[:n-1]

								sink.Next(newContext.Buffer)

								break
							}
						}

						critical.Leave(&x.Section)
					}
				}

				closingNotifier := obs.ClosingSelector(t.Value)
				closingNotifier.Subscribe(ctx, observer.Sink)

			case t.HasError:
				critical.Close(&x.Section)

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
					c.Buffer = append(c.Buffer, t.Value)
				}

				critical.Leave(&x.Section)

			case t.HasError:
				critical.Close(&x.Section)

				sink(t)

			default:
				critical.Close(&x.Section)

				for _, c := range x.Contexts {
					if ctx.Err() != nil {
						return
					}
					c.Cancel()
					sink.Next(c.Buffer)
				}

				sink(t)
			}
		}
	})
}
