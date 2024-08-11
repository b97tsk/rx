// Package rx is a reactive programming library for Go, inspired by
// https://reactivex.io/ (mostly RxJS).
//
// # Observers
//
//	type Observer[T any] func(n Notification[T])
//
// An [Observer] is a consumer of notifications delivered by an [Observable].
//
// An [Observer] is usually created and passed to [Observable.Subscribe] method
// when subscribing to an [Observable].
//
// # Observables
//
//	type Observable[T any] func(c Context, o Observer[T])
//
// An [Observable] is a collection of future values, waiting to become a flow
// of data. Subscribing an [Observer] to an [Observable] makes it happen.
// When an [Observable] is subscribed, its values, when available, are emitted
// to a given [Observer].
//
// There are four kinds of notifications:
//   - [Next]: value notifications.
//     An [Observable] can emit zero or more value notifications before
//     emitting one of the other kinds of notifications as a termination.
//   - [Complete]: completion notifications.
//     After emitting some value notifications, an [Observable] can emit one
//     completion notification as a termination.
//   - [Error]: error notifications.
//     After emitting some value notifications, an [Observable] can emit one
//     error notification as a termination.
//   - [Stop]: stop notifications.
//     After emitting some value notifications, an [Observable] can emit one
//     stop notification as a termination.
//
// Both [Error] notifications and [Stop] notifications carry an error value.
// The main differences between them are as follows:
//   - [Error] notifications can be caught and replaced, whereas
//     [Stop] notifications must not;
//   - [Error] notifications can be as deterministic as [Next] or [Complete]
//     notifications, whereas [Stop] notifications are not (continue reading
//     for more about this).
//
// # Operators
//
//	type Operator[T, R any] interface {
//		Apply(source Observable[T]) Observable[R]
//	}
//
// An [Operator] is an operation on an [Observable]. When applied, they do not
// change the existing [Observable] value. Instead, they return a new one,
// whose subscription logic is based on the first [Observable].
//
// There are many kinds of Operators in this library.
// Here is a list of what Operators can do:
//   - Filtering: [Filter], [First], [Last], [Sample], [Skip], [Take], etc.
//   - Transforming: [Map], [GroupBy], [Pairwise], [Scan], etc.
//   - Combining: [ConcatWith], [MergeWith], [StartWith], etc.
//   - Aggregating: [Reduce], [ToSlice], etc.
//   - Error Handling: [Catch], [Retry], etc.
//   - Backpressure: [OnBackpressureBuffer], [OnBackpressureCongest], etc.
//   - Utility: [Delay], [Do], [Materialize], [Timeout], etc.
//
// Previously, [Operator] was also a function type like [Observable] and
// [Observer]. It was changed to be an interface type for one reason:
// implementations can carry additional methods for setting extra options.
// For example, [MergeMap] has two extra options:
// [MergeMapOperator.WithBuffering] and [MergeMapOperator.WithConcurrency],
// and this is how they are specified when using a [MergeMap]:
//
//	op := MergeMap(f).WithBuffering().WithConcurrency(3)
//
// # Chaining Operators
//
// To chain multiple Operators, do either this:
//
//	ob1 := op1.Apply(source)
//	ob2 := op2.Apply(ob1)
//	ob3 := op3.Apply(ob2)
//
// or this:
//
//	ob := Pipe3(source, op1, op2, op3)
//
// There are 9 Pipe functions in this library, from [Pipe1] to [Pipe9].
// For different number of Operators, use different Pipe function.
//
// When there are really too many Operators to chain, do either this:
//
//	ob1 := Pipe5(source, op1, op2, op3, op4, op5)
//	ob2 := Pipe5(ob1, op6, op7, op8, op9, op10)
//	ob3 := Pipe5(ob2, op11, op12, op13, op14, op15)
//
// or this:
//
//	ob := Pipe3(source,
//		Compose5(op1, op2, op3, op4, op5),
//		Compose5(op6, op7, op8, op9, op10),
//		Compose5(op11, op12, op13, op14, op15),
//	)
//
// There are 8 Compose functions in this library, from [Compose2] to
// [Compose9].
//
// # Concurrency Safety
//
// Notifications emitted by an [Observable] may come from any started
// goroutine, but they are guaranteed to be in sequence, one after another.
//
// Operators in a chain may run in different goroutines.
// In the following code:
//
//	Pipe3(ob, op1, op2, op3).Subscribe(c, o)
//
// Race conditions could happen for any two of ob, op1, op2, op3 and o.
//
// Race conditions could also happen for any two Observables, however, not
// every [Operator] or [Observable] has concurrent behavior.
//
// The following operations may cause concurrent behavior:
//   - [Context] cancellation (due to use of [Context.AfterFunc]);
//   - [AuditTime], [DebounceTime] and [ThrottleTime] (due to use of [Timer]);
//   - [Channelize] (due to use of [Context.Go]);
//   - [Delay] (due to use of [Timer]);
//   - [Go] (due to use of [Context.Go]);
//   - OnBackpressure operators (due to use of [Channelize]);
//   - [SampleTime] (due to use of [Ticker]);
//   - [Ticker] and [Timer] (due to use of [Context.Go]);
//   - [Timeout] (due to use of [time.AfterFunc]);
//   - [Zip2] to [Zip9] (due to use of [Context.Go]).
//
// The following operations may cause concurrent behavior due to [Context]
// cancellation:
//   - CombineLatest operators (due to use of [Synchronize]);
//   - [Connect] (due to use of [Multicast]);
//   - [GroupBy] (due to use of [Multicast]);
//   - Merge operators (due to use of [Synchronize]);
//   - [Multicast] and other relatives (due to use of [Context.AfterFunc] and
//     [Synchronize]);
//   - [Never] (due to use of [Context.AfterFunc]);
//   - [Share] (due to use of [Context.AfterFunc] and [Multicast]);
//   - [Synchronize] (due to use of [Context.AfterFunc]);
//   - [Unicast] and other relatives (due to use of [Context.AfterFunc]);
//   - WithLatestFrom operators (due to use of [Synchronize]);
//   - ZipWithBuffering operators (due to use of [Synchronize]).
//
// Since [Context] cancellations are very common in this library, and that
// a [Context] cancellation usually results in a [Stop] notification, emitted
// in a goroutine started by [Context.AfterFunc] or [Context.Go], handling
// [Stop] notifications must take extra precaution. The problem is that,
// [Stop] notifications are not deterministic. They may just come from random
// goroutines. If that happens, one would have to deal with race conditions.
//
// It's very common that an [Observable], when subscribed, also subscribes to
// other Observables.
// In this library, inner Observables are usually subscribed in the same
// goroutine where the outer one is being subscribed. However,
//   - Observables created by [Concat], [ConcatWith] or Concat(All|Map|MapTo)
//     with source buffering on, may or may not subscribe to inner Observables
//     (excluding the first one) in separate goroutines, depending on whether
//     inner Observables cause concurrent behavior;
//   - Observables created by [Go] always subscribe to their source
//     Observables in their own goroutines;
//   - Observables created by Merge(All|Map|MapTo) with source buffering on,
//     may or may not subscribe to inner Observables in separate goroutines,
//     depending on whether inner Observables cause concurrent behavior;
//   - Observables created by Zip[2-9] always subscribe to inner Observables
//     in their own goroutines.
//
// When in doubt, read the code.
package rx
