package rx

import (
	"context"
	"time"
)

// A BufferTimeConfigure is a configure for BufferTime.
type BufferTimeConfigure struct {
	TimeSpan         time.Duration
	CreationInterval time.Duration
	MaxBufferSize    int
}

// Use creates an Operator from this configure.
func (configure BufferTimeConfigure) Use() Operator {
	return func(source Observable) Observable {
		obs := bufferTimeObservable{source, configure}
		return Create(obs.Subscribe)
	}
}

type bufferTimeObservable struct {
	Source Observable
	BufferTimeConfigure
}

type bufferTimeContext struct {
	Cancel context.CancelFunc
	Buffer []interface{}
}

func (obs bufferTimeObservable) Subscribe(ctx context.Context, sink Observer) {
	type X struct {
		Contexts []*bufferTimeContext
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	var closeContext func(*bufferTimeContext)

	openContextLocked := func(x *X) {
		ctx, cancel := context.WithCancel(ctx)
		newContext := &bufferTimeContext{Cancel: cancel}
		x.Contexts = append(x.Contexts, newContext)
		scheduleOnce(ctx, obs.TimeSpan, func() {
			closeContext(newContext)
		})
	}

	openContext := func() {
		if x, ok := <-cx; ok {
			openContextLocked(x)
			cx <- x
		}
	}

	closeContext = func(toBeClosed *bufferTimeContext) {
		toBeClosed.Cancel()
		if x, ok := <-cx; ok {
			for i, c := range x.Contexts {
				if c == toBeClosed {
					copy(x.Contexts[i:], x.Contexts[i+1:])
					n := len(x.Contexts)
					x.Contexts[n-1] = nil
					x.Contexts = x.Contexts[:n-1]
					sink.Next(toBeClosed.Buffer)
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
				var bufferFullContexts []*bufferTimeContext
				for _, c := range x.Contexts {
					c.Buffer = append(c.Buffer, t.Value)
					if len(c.Buffer) == obs.MaxBufferSize {
						bufferFullContexts = append(bufferFullContexts, c)
					}
				}

				cx <- x

				for _, c := range bufferFullContexts {
					closeContext(c)
				}

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

// BufferTime buffers the source Observable values for a specific time period.
//
// BufferTime collects values from the past as a slice, and emits those slices
// periodically in time.
func (Operators) BufferTime(timeSpan time.Duration) Operator {
	return BufferTimeConfigure{TimeSpan: timeSpan}.Use()
}
