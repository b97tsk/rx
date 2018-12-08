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

	var (
		buffer []interface{}
		try    cancellableLocker
	)

	op.Notifier.Subscribe(ctx, func(t Notification) {
		if try.Lock() {
			switch {
			case t.HasValue:
				value := buffer
				buffer = nil
				sink.Next(value)
				try.Unlock()
			default:
				try.CancelAndUnlock()
				sink(t)
			}
		}
	})

	select {
	case <-ctx.Done():
		return canceledCtx, nothingToDo
	default:
	}

	source.Subscribe(ctx, func(t Notification) {
		if try.Lock() {
			switch {
			case t.HasValue:
				buffer = append(buffer, t.Value)
				try.Unlock()
			default:
				try.CancelAndUnlock()
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
