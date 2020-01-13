package rx

import (
	"context"
)

type catchObservable struct {
	Source   Observable
	Selector func(error) Observable
}

func (obs catchObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	obs.Source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			sink(t)
		case t.HasError:
			obs := obs.Selector(t.Error)
			obs.Subscribe(ctx, sink)
		default:
			sink(t)
		}
	})

	return ctx, cancel
}

// Catch catches errors on the Observable to be handled by returning a new
// Observable.
func (Operators) Catch(selector func(error) Observable) Operator {
	return func(source Observable) Observable {
		return catchObservable{source, selector}.Subscribe
	}
}
