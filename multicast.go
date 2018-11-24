package rx

import (
	"context"
)

// Multicast creates an Observable that calls subjectFactory function to create
// a Subject, passes this Subject to selector function, where you have a chance
// to make more subscriptions to that Subject, then subscribes the source
// Observable to emit values to that Subject. Moreover, the selector function
// returns an Observable which is subscribed to the sink Observer; if it
// completes or emits an error, all subscriptions must be canceled.
func (Operators) Multicast(subjectFactory func() *Subject, selector func(context.Context, *Subject) Observable) OperatorFunc {
	return func(source Observable) Observable {
		return source.Lift(
			func(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(ctx)
				subject := subjectFactory()
				obs := selector(ctx, subject)
				obs.Subscribe(ctx, Finally(sink, cancel))
				select {
				case <-ctx.Done():
					return canceledCtx, nothingToDo
				default:
				}
				source.Subscribe(ctx, subject.Observer)
				return ctx, cancel
			},
		)
	}
}
