package rx

import (
	"context"

	"github.com/b97tsk/rx/x/atomic"
)

type skipUntilObservable struct {
	Source   Observable
	Notifier Observable
}

func (obs skipUntilObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
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

		obs.Notifier.Subscribe(ctx, observer.Notify)
	}

	if ctx.Err() != nil {
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

		obs.Source.Subscribe(ctx, observer.Notify)
	}

	return ctx, cancel
}

// SkipUntil creates an Observable that skips items emitted by the source
// Observable until a second Observable emits an item.
func (Operators) SkipUntil(notifier Observable) Operator {
	return func(source Observable) Observable {
		return skipUntilObservable{source, notifier}.Subscribe
	}
}
