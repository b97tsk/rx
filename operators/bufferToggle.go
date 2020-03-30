package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

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
	type X struct {
		Contexts []*bufferToggleContext
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	obs.Openings.Subscribe(ctx, func(t rx.Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				ctx, cancel := context.WithCancel(ctx)
				newContext := &bufferToggleContext{Cancel: cancel}
				x.Contexts = append(x.Contexts, newContext)

				cx <- x

				var observer rx.Observer
				observer = func(t rx.Notification) {
					observer = rx.NopObserver
					cancel()
					if x, ok := <-cx; ok {
						if t.HasError {
							close(cx)
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
						cx <- x
					}
				}

				closingNotifier := obs.ClosingSelector(t.Value)
				closingNotifier.Subscribe(ctx, observer.Notify)

			case t.HasError:
				close(cx)
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
					c.Buffer = append(c.Buffer, t.Value)
				}

				cx <- x

			case t.HasError:
				close(cx)
				sink(t)

			default:
				close(cx)
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

// BufferToggle buffers the source Observable values starting from an emission
// from openings and ending when the output of closingSelector emits.
//
// BufferToggle collects values from the past as a slice, starts collecting
// only when opening emits, and calls the closingSelector function to get an
// Observable that tells when to close the buffer.
func BufferToggle(openings rx.Observable, closingSelector func(interface{}) rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := bufferToggleObservable{source, openings, closingSelector}
		return rx.Create(obs.Subscribe)
	}
}
