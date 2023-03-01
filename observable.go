package rx

import (
	"context"
)

// An Observable is a collection of future values.
// When an Observable is subscribed, its values, when available,
// are emitted to a given [Observer].
//
// To cancel a subscription to an Observable, cancel the given [context.Context]
// value.
//
// An Observable must honor the Observable protocol, that is,
// no more emissions can pass to the given Observer after
// a notification of error or completion has passed to it.
// Violations of this protocol result in undefined behaviors.
//
// An Observable must honor the cancellation of the given context.
// When the cancellation of the given context is detected, an Observable must
// emits a notification of error to the given Observer as soon as possible.
//
// Observables are expected to be sequential. If you want to do something
// parallel, you will need to divide it (as an Observable) into pieces
// (Observables), process (subscribe to) them concurrently, and then later
// or in the meantime, merge (flatten) them together (back into one sequential
// Observable). Generally, these can be done by one single [Operator],
// for example, [MergeMap].
type Observable[T any] func(ctx context.Context, sink Observer[T])

// Subscribe invokes an execution of an Observable.
//
// Subscribing to a nil Observable results in an error notification of ErrNil.
func (obs Observable[T]) Subscribe(ctx context.Context, sink Observer[T]) {
	if obs == nil {
		sink.Error(ErrNil)
		return
	}

	obs(ctx, sink)
}

// NewObservable creates an Observable from f.
func NewObservable[T any](f func(ctx context.Context, sink Observer[T])) Observable[T] {
	return f
}
