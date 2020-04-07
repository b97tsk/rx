package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/x/atomic"
)

type skipUntilObservable struct {
	Source   rx.Observable
	Notifier rx.Observable
}

func (obs skipUntilObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	sinkNoMutex := sink
	sink = rx.Mutex(sinkNoMutex)

	var noSkipping atomic.Uint32

	{
		ctx, cancel := context.WithCancel(ctx)

		var observer rx.Observer

		observer = func(t rx.Notification) {
			switch {
			case t.HasValue:
				noSkipping.Store(1)
				observer = rx.Noop
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
		return
	}

	{
		var observer rx.Observer

		observer = func(t rx.Notification) {
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
}

// SkipUntil creates an Observable that skips items emitted by the source
// Observable until a second Observable emits an item.
func SkipUntil(notifier rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := skipUntilObservable{source, notifier}
		return rx.Create(obs.Subscribe)
	}
}
