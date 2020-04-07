package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

type isEmptyObservable struct {
	Source rx.Observable
}

func (obs isEmptyObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	var observer rx.Observer
	observer = func(t rx.Notification) {
		switch {
		case t.HasValue:
			observer = rx.Noop
			sink.Next(false)
			sink.Complete()
		case t.HasError:
			sink(t)
		default:
			sink.Next(true)
			sink.Complete()
		}
	}
	obs.Source.Subscribe(ctx, observer.Notify)
}

// IsEmpty creates an Observable that emits true if the source Observable
// emits no items, otherwise, it emits false.
func IsEmpty() rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := isEmptyObservable{source}
		return rx.Create(obs.Subscribe)
	}
}
