package operators

import (
	"context"
	"time"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/critical"
)

// BufferTime buffers the source Observable values for a specific time period.
//
// BufferTime collects values from the past as a slice, and emits those slices
// periodically in time.
func BufferTime(d time.Duration) rx.Operator {
	return BufferTimeConfigure{TimeSpan: d}.Make()
}

// A BufferTimeConfigure is a configure for BufferTime.
type BufferTimeConfigure struct {
	TimeSpan         time.Duration
	CreationInterval time.Duration
	MaxBufferSize    int
}

// Make creates an Operator from this configure.
func (configure BufferTimeConfigure) Make() rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return bufferTimeObservable{source, configure}.Subscribe
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
	ctx, cancel := context.WithCancel(ctx)
	sink = sink.WithCancel(cancel)

	var x struct {
		critical.Section
		Contexts []*bufferTimeContext
	}

	var closeContext func(*bufferTimeContext)

	obsTimer := rx.Timer(obs.TimeSpan)
	openContextLocked := func() {
		ctx, cancel := context.WithCancel(ctx)
		newContext := &bufferTimeContext{Cancel: cancel}
		x.Contexts = append(x.Contexts, newContext)
		obsTimer.Subscribe(ctx, func(t rx.Notification) {
			if t.HasValue {
				return
			}
			closeContext(newContext)
		})
	}

	openContext := func() {
		if critical.Enter(&x.Section) {
			openContextLocked()
			critical.Leave(&x.Section)
		}
	}

	closeContext = func(toBeClosed *bufferTimeContext) {
		toBeClosed.Cancel()
		if critical.Enter(&x.Section) {
			for i, c := range x.Contexts {
				if c == toBeClosed {
					copy(x.Contexts[i:], x.Contexts[i+1:])
					n := len(x.Contexts)
					x.Contexts[n-1] = nil
					x.Contexts = x.Contexts[:n-1]
					sink.Next(toBeClosed.Buffer)
					if obs.CreationInterval <= 0 {
						openContextLocked()
					}
					break
				}
			}
			critical.Leave(&x.Section)
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
		if critical.Enter(&x.Section) {
			switch {
			case t.HasValue:
				var bufferFullContexts []*bufferTimeContext
				for _, c := range x.Contexts {
					c.Buffer = append(c.Buffer, t.Value)
					if len(c.Buffer) == obs.MaxBufferSize {
						bufferFullContexts = append(bufferFullContexts, c)
					}
				}

				critical.Leave(&x.Section)

				for _, c := range bufferFullContexts {
					closeContext(c)
				}

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
