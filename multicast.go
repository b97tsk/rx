package rx

import (
	"context"
)

// Multicast creates an Observable that calls the subjectFactory function
// to create a Subject, subscribes the Subject to emit values to the sink
// Observer, then subscribes the source Observable to emit values to the
// Subject. In the subjectFactory function, you can create subscriptions to
// the Subject as many as you need, just before it returns.
func (Operators) Multicast(subjectFactory func(context.Context) Subject) OperatorFunc {
	return func(source Observable) Observable {
		return source.Lift(
			func(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(ctx)
				subject := subjectFactory(ctx)
				subject.Subscribe(ctx, Finally(sink, cancel))
				if isDone(ctx) {
					return canceledCtx, nothingToDo
				}
				source.Subscribe(ctx, subject.Observer)
				return ctx, cancel
			},
		)
	}
}
