package rx

import (
	"context"
)

type bufferWhenOperator struct {
	ClosingSelector func() Observable
}

func (op bufferWhenOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	var buffer struct {
		cancellableLocker
		List []interface{}
	}

	var (
		openBuffer     func()
		avoidRecursive avoidRecursiveCalls
	)

	openBuffer = func() {
		if isDone(ctx) {
			return
		}

		ctx, cancel := context.WithCancel(ctx)

		var observer Observer
		observer = func(t Notification) {
			observer = NopObserver
			cancel()
			if buffer.Lock() {
				if t.HasError {
					buffer.CancelAndUnlock()
					sink(t)
					return
				}
				sink.Next(buffer.List)
				buffer.List = nil
				buffer.Unlock()
				avoidRecursive.Do(openBuffer)
			}
		}

		closingNotifier := op.ClosingSelector()
		closingNotifier.Subscribe(ctx, observer.Notify)
	}

	avoidRecursive.Do(openBuffer)

	if isDone(ctx) {
		return Done()
	}

	source.Subscribe(ctx, func(t Notification) {
		if buffer.Lock() {
			switch {
			case t.HasValue:
				buffer.List = append(buffer.List, t.Value)
				buffer.Unlock()
			default:
				buffer.CancelAndUnlock()
				sink.Next(buffer.List)
				sink(t)
			}
		}
	})

	return ctx, cancel
}

// BufferWhen buffers the source Observable values, using a factory function
// of closing Observables to determine when to close, emit, and reset the
// buffer.
//
// BufferWhen collects values from the past as a slice. When it starts
// collecting values, it calls a function that returns an Observable that
// tells when to close the buffer and restart collecting.
//
// Dead loop could happen if closing Observables emit a value or complete as
// soon as they are subscribed to.
func (Operators) BufferWhen(closingSelector func() Observable) OperatorFunc {
	return func(source Observable) Observable {
		op := bufferWhenOperator{closingSelector}
		return source.Lift(op.Call)
	}
}
