package rx

import (
	"context"
	"sync/atomic"
)

type skipUntilOperator struct {
	source   Operator
	notifier Observable
}

func (op skipUntilOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()
	noSkipping := uint32(0)
	hasCompleted := uint32(0)
	ob = Normalize(ob)

	op.notifier.Subscribe(ctx, ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			atomic.StoreUint32(&noSkipping, 1)
		case t.HasError:
			ob.Error(t.Value.(error))
			cancel()
		default:
			if atomic.CompareAndSwapUint32(&hasCompleted, 0, 1) {
				break
			}
			ob.Complete()
			cancel()
		}
	}))

	select {
	case <-done:
		return ctx, cancel
	default:
	}

	op.source.Call(ctx, ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			if atomic.LoadUint32(&noSkipping) != 0 {
				ob.Next(t.Value)
			}
		case t.HasError:
			ob.Error(t.Value.(error))
			cancel()
		default:
			if atomic.CompareAndSwapUint32(&hasCompleted, 0, 1) {
				break
			}
			ob.Complete()
			cancel()
		}
	}))

	return ctx, cancel
}

// SkipUntil creates an Observable that skips items emitted by the source
// Observable until a second Observable emits an item.
func (o Observable) SkipUntil(notifier Observable) Observable {
	op := skipUntilOperator{
		source:   o.Op,
		notifier: notifier,
	}
	return Observable{op}
}
