package rx

import (
	"context"
)

type bufferOperator struct {
	Notifier Observable
}

func (op bufferOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()

	var (
		buffer []interface{}
		try    cancellableLocker
	)

	op.Notifier.Subscribe(ctx, func(t Notification) {
		if try.Lock() {
			switch {
			case t.HasValue:
				defer try.Unlock()
				value := buffer
				buffer = nil
				sink.Next(value)
			default:
				try.CancelAndUnlock()
				sink(t)
				cancel()
			}
		}
	})

	select {
	case <-done:
		return ctx, cancel
	default:
	}

	source.Subscribe(ctx, func(t Notification) {
		if try.Lock() {
			switch {
			case t.HasValue:
				defer try.Unlock()
				buffer = append(buffer, t.Value)
			default:
				try.CancelAndUnlock()
				sink(t)
				cancel()
			}
		}
	})

	return ctx, cancel
}

// Buffer buffers the source Observable values until notifier emits.
//
// Buffer collects values from the past as a slice, and emits that slice
// only when another Observable emits.
func (o Observable) Buffer(notifier Observable) Observable {
	op := bufferOperator{notifier}
	return o.Lift(op.Call)
}
