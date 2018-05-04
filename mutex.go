package rx

import (
	"context"
)

type mutexOperator struct {
	source Operator
}

func (op mutexOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	try := cancellableLocker{}
	return op.source.Call(ctx, ObserverFunc(func(t Notification) {
		if try.Lock() {
			switch {
			case t.HasValue:
				ob.Next(t.Value)
				try.Unlock()
			case t.HasError:
				try.CancelAndUnlock()
				ob.Error(t.Value.(error))
			default:
				try.CancelAndUnlock()
				ob.Complete()
			}
		}
	}))
}

// Mutex creates an Observable that mirrors the source Observable in a mutually
// exclusive way.
func (o Observable) Mutex() Observable {
	op := mutexOperator{o.Op}
	return Observable{op}
}
