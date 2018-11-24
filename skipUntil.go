package rx

import (
	"context"
	"sync/atomic"
)

type skipUntilOperator struct {
	Notifier Observable
}

func (op skipUntilOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sinkNoMutex := Finally(sink, cancel)
	sink = Mutex(sinkNoMutex)

	var noSkipping uint32

	{
		ctx, cancel := context.WithCancel(ctx)

		var observer Observer

		observer = func(t Notification) {
			switch {
			case t.HasValue:
				atomic.StoreUint32(&noSkipping, 1)
				observer = NopObserver
				cancel()
			case t.HasError:
				sink(t)
			default:
				// do nothing
			}
		}

		op.Notifier.Subscribe(ctx, observer.Notify)
	}

	select {
	case <-ctx.Done():
		return canceledCtx, nothingToDo
	default:
	}

	{
		var observer Observer

		observer = func(t Notification) {
			switch {
			case t.HasValue:
				if atomic.LoadUint32(&noSkipping) != 0 {
					observer = sinkNoMutex
					sinkNoMutex(t)
				}
			default:
				sink(t)
			}
		}

		source.Subscribe(ctx, observer.Notify)
	}

	return ctx, cancel
}

// SkipUntil creates an Observable that skips items emitted by the source
// Observable until a second Observable emits an item.
func (Operators) SkipUntil(notifier Observable) OperatorFunc {
	return func(source Observable) Observable {
		op := skipUntilOperator{notifier}
		return source.Lift(op.Call)
	}
}
