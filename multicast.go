package rx

import (
	"context"
)

type multicastOperator struct {
	subjectFactory func() *Subject
	selector       func(context.Context, *Subject) Observable
}

func (op multicastOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	subject := op.subjectFactory()
	obsv := op.selector(ctx, subject)
	obsv.Subscribe(ctx, Finally(sink, cancel))
	select {
	case <-ctx.Done():
	default:
		source.Subscribe(ctx, subject.Observer)
	}
	return ctx, cancel
}

// Multicast creates an Observable that calls subjectFactory function to create
// a Subject, passes this Subject to selector function, where you have a chance
// to make more subscriptions to that Subject, then subscribes the source
// Observable to emit values to that Subject. Moreover, the selector function
// returns an Observable which is subscribed to the Observer, if it completes
// or emits an error, all subscriptions shall be canceled.
func (o Observable) Multicast(subjectFactory func() *Subject, selector func(context.Context, *Subject) Observable) Observable {
	op := multicastOperator{subjectFactory, selector}
	return o.Lift(op.Call)
}
