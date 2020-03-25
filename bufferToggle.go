package rx

import (
	"context"
)

type bufferToggleObservable struct {
	Source          Observable
	Openings        Observable
	ClosingSelector func(interface{}) Observable
}

type bufferToggleContext struct {
	Cancel context.CancelFunc
	Buffer []interface{}
}

func (obs bufferToggleObservable) Subscribe(parent context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx := NewContext(parent)

	sink = DoAtLast(sink, ctx.AtLast)

	type X struct {
		Contexts []*bufferToggleContext
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	obs.Openings.Subscribe(ctx, func(t Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				ctx, cancel := context.WithCancel(ctx)
				newContext := &bufferToggleContext{Cancel: cancel}
				x.Contexts = append(x.Contexts, newContext)

				cx <- x

				var observer Observer
				observer = func(t Notification) {
					observer = NopObserver
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
		return ctx, ctx.Cancel
	}

	obs.Source.Subscribe(ctx, func(t Notification) {
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

	return ctx, ctx.Cancel
}

// BufferToggle buffers the source Observable values starting from an emission
// from openings and ending when the output of closingSelector emits.
//
// BufferToggle collects values from the past as a slice, starts collecting
// only when opening emits, and calls the closingSelector function to get an
// Observable that tells when to close the buffer.
func (Operators) BufferToggle(openings Observable, closingSelector func(interface{}) Observable) Operator {
	return func(source Observable) Observable {
		return bufferToggleObservable{source, openings, closingSelector}.Subscribe
	}
}
