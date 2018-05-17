package rx

import (
	"context"
)

type mutexOperator struct{}

func (op mutexOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	var try cancellableLocker
	return source.Subscribe(ctx, func(t Notification) {
		if try.Lock() {
			switch {
			case t.HasValue:
				defer try.Unlock()
				ob.Next(t.Value)
			case t.HasError:
				try.CancelAndUnlock()
				ob.Error(t.Value.(error))
			default:
				try.CancelAndUnlock()
				ob.Complete()
			}
		}
	})
}

// Mutex creates an Observable that mirrors the source Observable in a mutually
// exclusive way.
func (o Observable) Mutex() Observable {
	op := mutexOperator{}
	return o.Lift(op.Call)
}
