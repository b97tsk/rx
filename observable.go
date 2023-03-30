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
// An Observable must honor the Observable protocol:
//   - An Observable can emit zero or more values of specific type;
//   - An Observable must emit a notification of error or completion as
//     a termination;
//   - An Observable must not emit anything after a termination.
//
// Missing a termination may lead to memory leaks;
// Emissions after a termination cause undefined behavior.
//
// An Observable must honor the cancellation of the given context.
// When the cancellation of the given context is detected, an Observable must
// emit a notification of error (as a termination) to the given Observer
// as soon as possible.
//
// An Observable must use [Go] function rather than go statements to start
// new goroutines; otherwise, your program might panic randomly when using
// any of the Blocking methods.
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
