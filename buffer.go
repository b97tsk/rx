package rx

import (
	"context"
)

type bufferOperator struct {
	Notifier Observable
}

func (op bufferOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	var buffer struct {
		cancellableLocker
		List []interface{}
	}

	op.Notifier.Subscribe(ctx, func(t Notification) {
		if buffer.Lock() {
			switch {
			case t.HasValue:
				sink.Next(buffer.List)
				buffer.List = nil
				buffer.Unlock()
			default:
				buffer.CancelAndUnlock()
				sink(t)
			}
		}
	})

	if isDone(ctx) {
		return canceledCtx, nothingToDo
	}

	source.Subscribe(ctx, func(t Notification) {
		if buffer.Lock() {
			switch {
			case t.HasValue:
				buffer.List = append(buffer.List, t.Value)
				buffer.Unlock()
			default:
				buffer.CancelAndUnlock()
				sink(t)
			}
		}
	})

	return ctx, cancel
}

// Buffer buffers the source Observable values until notifier emits.
//
// Buffer collects values from the past as a slice, and emits that slice
// only when another Observable emits.
func (Operators) Buffer(notifier Observable) OperatorFunc {
	return func(source Observable) Observable {
		op := bufferOperator{notifier}
		return source.Lift(op.Call)
	}
}
