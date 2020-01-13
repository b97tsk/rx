package rx

import (
	"context"
)

type dematerializeObservable struct {
	Source Observable
}

func (obs dematerializeObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	var observer Observer
	observer = func(t Notification) {
		switch {
		case t.HasValue:
			if t, ok := t.Value.(Notification); ok {
				switch {
				case t.HasValue:
					sink(t)
				default:
					observer = NopObserver
					sink(t)
				}
			} else {
				observer = NopObserver
				sink.Error(ErrNotNotification)
			}
		default:
			sink(t)
		}
	}
	obs.Source.Subscribe(ctx, observer.Notify)

	return ctx, cancel
}

// Dematerialize converts an Observable of Notification objects into the
// emissions that they represent. It's the opposite of Materialize.
func (Operators) Dematerialize() Operator {
	return func(source Observable) Observable {
		return dematerializeObservable{source}.Subscribe
	}
}
