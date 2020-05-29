package operators

import (
	"context"
	"time"

	"github.com/b97tsk/rx"
)

// A BufferTimeConfigure is a configure for BufferTime.
type BufferTimeConfigure struct {
	TimeSpan         time.Duration
	CreationInterval time.Duration
	MaxBufferSize    int
}

// Use creates an Operator from this configure.
func (configure BufferTimeConfigure) Use() rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := bufferTimeObservable{source, configure}
		return rx.Create(obs.Subscribe)
	}
}

type bufferTimeObservable struct {
	Source rx.Observable
	BufferTimeConfigure
}

type bufferTimeContext struct {
	Cancel context.CancelFunc
	Buffer []interface{}
}

func (obs bufferTimeObservable) Subscribe(ctx context.Context, sink rx.Observer) {
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
		rx.Timer(obs.TimeSpan).Subscribe(ctx, func(t rx.Notification) {
			if t.HasValue {
				return
			}
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
		rx.Ticker(obs.CreationInterval).Subscribe(ctx, func(t rx.Notification) {
			if t.HasValue {
				openContext()
			}
		})
	}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
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
func BufferTime(timeSpan time.Duration) rx.Operator {
	return BufferTimeConfigure{TimeSpan: timeSpan}.Use()
}
