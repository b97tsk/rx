package rx

import (
	"context"

	"github.com/b97tsk/rx/x/atomic"
)

type skipUntilOperator struct {
	Notifier Observable
}

func (op skipUntilOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sinkNoMutex := Finally(sink, cancel)
	sink = Mutex(sinkNoMutex)

	var noSkipping atomic.Uint32

	{
		ctx, cancel := context.WithCancel(ctx)

		var observer Observer

		observer = func(t Notification) {
			switch {
			case t.HasValue:
				noSkipping.Store(1)
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

	if isDone(ctx) {
		return Done()
	}

	{
		var observer Observer

		observer = func(t Notification) {
			switch {
			case t.HasValue:
				if noSkipping.Equals(1) {
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
