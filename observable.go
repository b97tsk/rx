package rx

import (
	"context"
)

// An Observable is a collection of future values. When an Observable is
// subscribed, its values, when available, are emitted to the given Observer.
//
// To cancel a subscription to an Observable, cancel the given context.Context
// value.
//
// An Observable must honor the Observable contract, that is, no more
// emissions should pass to the given Observer after an error or a completion
// has passed to it. Violations of this contract result in undefined behaviors.
//
// An Observable should honor the cancellation of the given context. Either
// emits an error or a completion to the given Observer, when the cancellation
// of the given context is detected.
//
// Observables are expected to be sequential. If you want parallel, you
// must divide the work (as an Observable) you're going to do into pieces
// (Observables), process (subscribe to) them concurrently, and then later
// or in the meantime, merge (flatten) them together (back into one sequential
// Observable). Typically, these are done by one single Operator, for example,
// MergeMap.
type Observable[T any] func(ctx context.Context, sink Observer[T])

// Subscribe invokes an execution of an Observable.
func (obs Observable[T]) Subscribe(ctx context.Context, sink Observer[T]) {
	obs(ctx, sink)
}

// AsObservable converts f to an Observable.
func AsObservable[T any](f func(ctx context.Context, sink Observer[T])) Observable[T] {
	return f
}
