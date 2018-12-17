package rx

import (
	"context"
)

// Finally creates an Observer that passes all emissions to the specified
// Observer, in the case that an ERROR or COMPLETE emission is passed, makes
// a call to the specified function.
func Finally(sink Observer, finally func()) Observer {
	return func(t Notification) {
		if t.HasValue {
			sink(t)
			return
		}
		sink(t)
		finally()
	}
}

// Finally creates an Observable that mirrors the source Observable, in the
// case that an ERROR or COMPLETE emission is mirrored, makes a call to the
// specified function.
func (Operators) Finally(finally func()) OperatorFunc {
	return func(source Observable) Observable {
		return source.Lift(
			func(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
				return source.Subscribe(ctx, Finally(sink, finally))
			},
		)
	}
}
