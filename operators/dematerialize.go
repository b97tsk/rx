package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

type dematerializeObservable struct {
	Source rx.Observable
}

func (obs dematerializeObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	var observer rx.Observer
	observer = func(t rx.Notification) {
		switch {
		case t.HasValue:
			if t, ok := t.Value.(rx.Notification); ok {
				switch {
				case t.HasValue:
					sink(t)
				default:
					observer = rx.NopObserver
					sink(t)
				}
			} else {
				observer = rx.NopObserver
				sink.Error(rx.ErrNotNotification)
			}
		default:
			sink(t)
		}
	}
	obs.Source.Subscribe(ctx, observer.Notify)
}

// Dematerialize converts an Observable of Notification objects into the
// emissions that they represent. It's the opposite of Materialize.
func Dematerialize() rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := dematerializeObservable{source}
		return rx.Create(obs.Subscribe)
	}
}
