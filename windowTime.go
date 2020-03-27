package rx

import (
	"context"
	"time"
)

// A WindowTimeConfigure is a configure for WindowTime.
type WindowTimeConfigure struct {
	TimeSpan         time.Duration
	CreationInterval time.Duration
	MaxWindowSize    int
}

// Use creates an Operator from this configure.
func (configure WindowTimeConfigure) Use() Operator {
	return func(source Observable) Observable {
		obs := windowTimeObservable{source, configure}
		return Create(obs.Subscribe)
	}
}

type windowTimeObservable struct {
	Source Observable
	WindowTimeConfigure
}

type windowTimeContext struct {
	Cancel context.CancelFunc
	Window Subject
	Size   int
}

func (obs windowTimeObservable) Subscribe(ctx context.Context, sink Observer) {
	type X struct {
		Contexts []*windowTimeContext
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	var closeContext func(*windowTimeContext)

	openContextLocked := func(x *X) {
		ctx, cancel := context.WithCancel(ctx)
		newContext := &windowTimeContext{
			Cancel: cancel,
			Window: NewSubject(),
		}
		x.Contexts = append(x.Contexts, newContext)
		scheduleOnce(ctx, obs.TimeSpan, func() {
			closeContext(newContext)
		})
		sink.Next(newContext.Window.Observable)
	}

	openContext := func() {
		if x, ok := <-cx; ok {
			openContextLocked(x)
			cx <- x
		}
	}

	closeContext = func(toBeClosed *windowTimeContext) {
		toBeClosed.Cancel()
		if x, ok := <-cx; ok {
			for i, c := range x.Contexts {
				if c == toBeClosed {
					copy(x.Contexts[i:], x.Contexts[i+1:])
					n := len(x.Contexts)
					x.Contexts[n-1] = nil
					x.Contexts = x.Contexts[:n-1]
					toBeClosed.Window.Complete()
					if obs.CreationInterval <= 0 {
						openContextLocked(x)
					}
					break
				}
			}
			cx <- x
		}
	}

	openContext()

	if obs.CreationInterval > 0 {
		schedule(ctx, obs.CreationInterval, openContext)
	}

	obs.Source.Subscribe(ctx, func(t Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				var windowFullContexts []*windowTimeContext
				for _, c := range x.Contexts {
					c.Size++
					t.Observe(c.Window.Observer)
					if c.Size == obs.MaxWindowSize {
						windowFullContexts = append(windowFullContexts, c)
					}
				}

				cx <- x

				for _, c := range windowFullContexts {
					closeContext(c)
				}

			default:
				close(cx)
				for _, c := range x.Contexts {
					c.Cancel()
					t.Observe(c.Window.Observer)
				}
				sink(t)
			}
		}
	})
}

// WindowTime branches out the source Observable values as a nested Observable
// periodically in time.
//
// It's like BufferTime, but emits a nested Observable instead of a slice.
func (Operators) WindowTime(timeSpan time.Duration) Operator {
	return WindowTimeConfigure{TimeSpan: timeSpan}.Use()
}
