package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/atomic"
)

type skipUntilObservable struct {
	Source   rx.Observable
	Notifier rx.Observable
}

func (obs skipUntilObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	originalSink := sink
	sink = rx.Mutex(originalSink)

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

		obs.Notifier.Subscribe(ctx, observer.Sink)
	}

	if ctx.Err() != nil {
		return
	}

	if noSkipping.Equals(1) {
		obs.Source.Subscribe(ctx, originalSink)
		return
	}

	{
		var observer rx.Observer

		observer = func(t rx.Notification) {
			switch {
			case t.HasValue:
				if noSkipping.Equals(1) {
					observer = originalSink
					observer.Sink(t)
				}
			default:
				sink(t)
			}
		}

		obs.Source.Subscribe(ctx, observer.Sink)
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
