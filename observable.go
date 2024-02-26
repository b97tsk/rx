package rx

// An Observable is a collection of future values.
// When an Observable is subscribed, its values, when available,
// are emitted to a given [Observer].
//
// To cancel a subscription to an Observable, cancel the given [Context].
//
// An Observable must honor the Observable protocol:
//   - An Observable can emit zero or more values of specific type;
//   - An Observable must emit a notification of error or completion as
//     a termination;
//   - An Observable must not emit anything after a termination.
//
// Observables are expected to be sequential. Every emission emitted to
// the given [Observer] must be concurrency safe.
//
// An Observable must honor the cancellation of the given [Context].
// When the cancellation of the given [Context] is detected, an Observable must
// emit an error notification (as a termination) to the given Observer as soon
// as possible.
//
// If an Observable needs to start goroutines, it must use [Context.Go] to do
// so; if an Observable needs to start an asynchronous operation other than
// goroutines, it must call [Context.PreAsyncCall] to wrap what that Observable
// would do in that asynchronous operation, then call the function returned in
// that asynchronous operation instead.
//
// To achieve something in parallel, multiple Observables might be involved.
// There are a couple of functions and Operators in this library can handle
// multiple Observables, most of them can do things concurrently.
type Observable[T any] func(c Context, sink Observer[T])

// Subscribe invokes an execution of an Observable.
//
// Subscribing to a nil Observable results in an error notification of ErrNil.
func (obs Observable[T]) Subscribe(c Context, sink Observer[T]) {
	if obs == nil {
		sink.Error(ErrNil)
		return
	}

	obs(c, sink)
}

// NewObservable creates an Observable from f.
func NewObservable[T any](f func(c Context, sink Observer[T])) Observable[T] {
	return f
}
