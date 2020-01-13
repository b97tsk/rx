package rx

import (
	"context"
)

type countObservable struct {
	Source Observable
}

func (obs countObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	var count int
	return obs.Source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			count++
		case t.HasError:
			sink(t)
		default:
			sink.Next(count)
			sink.Complete()
		}
	})
}

// Count creates an Observable that counts the number of NEXT emissions on
// the source and emits that number when the source completes.
func (Operators) Count() OperatorFunc {
	return func(source Observable) Observable {
		return countObservable{source}.Subscribe
	}
}
