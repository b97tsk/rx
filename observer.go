package rx

// An Observer is a consumer of notifications delivered by an Observable.
type Observer[T any] func(n Notification[T])

// Next passes a value to sink.
func (sink Observer[T]) Next(v T) {
	sink(Next(v))
}

// Error passes an error to sink.
func (sink Observer[T]) Error(e error) {
	sink(Error[T](e))
}

// Complete passes a completion to sink.
func (sink Observer[T]) Complete() {
	sink(Complete[T]())
}

// Sink passes n to *sink.
//
// Sink also yields an Observer that is equivalent to:
//
//	func(n Notification[T]) { (*sink)(n) }
//
// Useful when you want to set *sink to another Observer at some point.
//
func (sink *Observer[T]) Sink(n Notification[T]) {
	(*sink)(n)
}

// Noop gives you an Observer that does nothing.
func Noop[T any](Notification[T]) {}

// AsObserver converts f to an Observer.
func AsObserver[T any](f func(n Notification[T])) Observer[T] { return f }
