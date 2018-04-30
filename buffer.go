package rx

import (
	"context"
)

type bufferOperator struct {
	source   Operator
	notifier Observable
}

func (op bufferOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()
	try := cancellableLocker{}
	buffer := []interface{}(nil)

	op.notifier.Subscribe(ctx, ObserverFunc(func(t Notification) {
		if try.Lock() {
			switch {
			case t.HasValue:
				value := buffer
				buffer = nil
				ob.Next(value)
				try.Unlock()
			case t.HasError:
				try.Cancel()
				ob.Error(t.Value.(error))
				cancel()
			default:
				try.Cancel()
				ob.Complete()
				cancel()
			}
		}
	}))

	select {
	case <-done:
		return ctx, cancel
	default:
	}

	op.source.Call(ctx, ObserverFunc(func(t Notification) {
		if try.Lock() {
			switch {
			case t.HasValue:
				buffer = append(buffer, t.Value)
				try.Unlock()
			case t.HasError:
				try.Cancel()
				ob.Error(t.Value.(error))
				cancel()
			default:
				try.Cancel()
				ob.Complete()
				cancel()
			}
		}
	}))

	return ctx, cancel
}

// Buffer buffers the source Observable values until notifier emits.
//
// Buffer collects values from the past as a slice, and emits that slice
// only when another Observable emits.
func (o Observable) Buffer(notifier Observable) Observable {
	op := bufferOperator{
		source:   o.Op,
		notifier: notifier,
	}
	return Observable{op}
}
