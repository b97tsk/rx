package rx

import (
	"context"
)

type bufferOperator struct {
	notifier Observable
}

func (op bufferOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()

	var (
		buffer []interface{}
		try    cancellableLocker
	)

	op.notifier.Subscribe(ctx, func(t Notification) {
		if try.Lock() {
			switch {
			case t.HasValue:
				defer try.Unlock()
				value := buffer
				buffer = nil
				ob.Next(value)
			case t.HasError:
				try.CancelAndUnlock()
				ob.Error(t.Value.(error))
				cancel()
			default:
				try.CancelAndUnlock()
				ob.Complete()
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
			case t.HasError:
				try.CancelAndUnlock()
				ob.Error(t.Value.(error))
				cancel()
			default:
				try.CancelAndUnlock()
				ob.Complete()
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
