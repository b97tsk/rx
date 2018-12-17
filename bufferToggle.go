package rx

import (
	"context"
)

type bufferToggleOperator struct {
	Openings        Observable
	ClosingSelector func(interface{}) Observable
}

type bufferToggleContext struct {
	Cancel context.CancelFunc
	Buffer []interface{}
}

func (op bufferToggleOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	var contexts struct {
		cancellableLocker
		List []*bufferToggleContext
	}

	op.Openings.Subscribe(ctx, func(t Notification) {
		if contexts.Lock() {
			switch {
			case t.HasValue:
				ctx, cancel := context.WithCancel(ctx)
				newContext := &bufferToggleContext{Cancel: cancel}
				contexts.List = append(contexts.List, newContext)
				contexts.Unlock()

				var observer Observer
				observer = func(t Notification) {
					observer = NopObserver
					cancel()
					if contexts.Lock() {
						if t.HasError {
							contexts.CancelAndUnlock()
							sink(t)
							return
						}
						for i, btc := range contexts.List {
							if btc == newContext {
								copy(contexts.List[i:], contexts.List[i+1:])
								n := len(contexts.List)
								contexts.List[n-1] = nil
								contexts.List = contexts.List[:n-1]
								sink.Next(newContext.Buffer)
								break
							}
						}
						contexts.Unlock()
					}
				}

				closingNotifier := op.ClosingSelector(t.Value)
				closingNotifier.Subscribe(ctx, observer.Notify)

			case t.HasError:
				contexts.CancelAndUnlock()
				sink(t)

			default:
				contexts.Unlock()
			}
		}
	})

	if isDone(ctx) {
		return Done()
	}

	source.Subscribe(ctx, func(t Notification) {
		if contexts.Lock() {
			switch {
			case t.HasValue:
				for _, btc := range contexts.List {
					btc.Buffer = append(btc.Buffer, t.Value)
				}
				contexts.Unlock()

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

// BufferToggle buffers the source Observable values starting from an emission
// from openings and ending when the output of closingSelector emits.
//
// BufferToggle collects values from the past as a slice, starts collecting
// only when opening emits, and calls the closingSelector function to get an
// Observable that tells when to close the buffer.
func (Operators) BufferToggle(openings Observable, closingSelector func(interface{}) Observable) OperatorFunc {
	return func(source Observable) Observable {
		op := bufferToggleOperator{openings, closingSelector}
		return source.Lift(op.Call)
	}
}
