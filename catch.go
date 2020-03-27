package rx

import (
	"context"
)

// Catch catches errors on the Observable to be handled by returning a new
// Observable.
func (Operators) Catch(selector func(error) Observable) Operator {
	return func(source Observable) Observable {
		return Create(
			func(ctx context.Context, sink Observer) {
				source.Subscribe(ctx, func(t Notification) {
					switch {
					case t.HasValue:
						sink(t)
					case t.HasError:
						obs := selector(t.Error)
						obs.Subscribe(ctx, sink)
					default:
						sink(t)
					}
				})
			},
		)
	}
}
