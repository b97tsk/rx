package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

type firstObservable struct {
	Source rx.Observable
}

func (obs firstObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	var observer rx.Observer
	observer = func(t rx.Notification) {
		switch {
		case t.HasValue:
			observer = rx.NopObserver
			sink(t)
			sink.Complete()
		case t.HasError:
			sink(t)
		default:
			sink.Error(rx.ErrEmpty)
		}
	}
	obs.Source.Subscribe(ctx, observer.Notify)
}

// First creates an Observable that emits only the first value (or the first
// value that meets some condition) emitted by the source Observable.
func First() rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := firstObservable{source}
		return rx.Create(obs.Subscribe)
	}
}
