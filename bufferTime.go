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

// MakeFunc creates an OperatorFunc from this type.
func (conf BufferTimeConfigure) MakeFunc() OperatorFunc {
	return MakeFunc(bufferTimeOperator(conf).Call)
}

type bufferTimeOperator BufferTimeConfigure

type bufferTimeContext struct {
	Cancel context.CancelFunc
	Buffer []interface{}
}

func (op bufferTimeOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	var contexts struct {
		cancellableLocker
		List []*bufferTimeContext
	}

	var closeContext func(*bufferTimeContext)

	openContextLocked := func() {
		ctx, cancel := context.WithCancel(ctx)
		newContext := &bufferTimeContext{Cancel: cancel}
		contexts.List = append(contexts.List, newContext)
		scheduleOnce(ctx, op.TimeSpan, func() {
			closeContext(newContext)
		})
	}

	openContext := func() {
		if contexts.Lock() {
			openContextLocked()
			contexts.Unlock()
		}
	}

	closeContext = func(toBeClosed *bufferTimeContext) {
		toBeClosed.Cancel()
		if contexts.Lock() {
			for i, btc := range contexts.List {
				if btc == toBeClosed {
					copy(contexts.List[i:], contexts.List[i+1:])
					n := len(contexts.List)
					contexts.List[n-1] = nil
					contexts.List = contexts.List[:n-1]
					sink.Next(toBeClosed.Buffer)
					if op.CreationInterval <= 0 {
						openContextLocked()
					}
					break
				}
			}
			contexts.Unlock()
		}
	}

	// No need to lock for the first time.
	openContextLocked()

	if op.CreationInterval > 0 {
		schedule(ctx, op.CreationInterval, openContext)
	}

	source.Subscribe(ctx, func(t Notification) {
		if contexts.Lock() {
			switch {
			case t.HasValue:
				var bufferFullContexts []*bufferTimeContext
				for _, btc := range contexts.List {
					btc.Buffer = append(btc.Buffer, t.Value)
					if len(btc.Buffer) == op.MaxBufferSize {
						bufferFullContexts = append(bufferFullContexts, btc)
					}
				}
				contexts.Unlock()

				for _, btc := range bufferFullContexts {
					closeContext(btc)
				}

			case t.HasError:
				contexts.CancelAndUnlock()
				sink(t)

			default:
				contexts.CancelAndUnlock()
				for _, btc := range contexts.List {
					if isDone(ctx) {
						return
					}
					btc.Cancel()
					sink.Next(btc.Buffer)
				}
				sink(t)
			}
		}
	})

	return ctx, cancel
}

// BufferTime buffers the source Observable values for a specific time period.
//
// BufferTime collects values from the past as a slice, and emits those slices
// periodically in time.
func (Operators) BufferTime(timeSpan time.Duration) OperatorFunc {
	return func(source Observable) Observable {
		op := bufferTimeOperator{TimeSpan: timeSpan}
		return source.Lift(op.Call)
	}
}
