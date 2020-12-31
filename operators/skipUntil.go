package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/atomic"
)

// SkipUntil skips items emitted by the source until a second Observable
// emits an item.
func SkipUntil(notifier rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return skipUntilObservable{source, notifier}.Subscribe
	}
}

type skipUntilObservable struct {
	Source   rx.Observable
	Notifier rx.Observable
}

func (obs skipUntilObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)

	originalSink := sink

	sink = sink.WithCancel(cancel).Mutex()

	var noSkipping atomic.Bool

	{
		ctx, cancel := context.WithCancel(ctx)

		var observer rx.Observer

		observer = func(t rx.Notification) {
			observer = rx.Noop

			cancel()

			switch {
			case t.HasValue:
				noSkipping.Store(true)
			case t.HasError:
				sink(t)
			}
		}

		obs.Notifier.Subscribe(ctx, observer.Sink)
	}

	if ctx.Err() != nil {
		return
	}

	if noSkipping.True() {
		obs.Source.Subscribe(ctx, originalSink)

		return
	}

	{
		var observer rx.Observer

		observer = func(t rx.Notification) {
			switch {
			case t.HasValue:
				if noSkipping.True() {
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
