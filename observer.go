package rx

// An Observer is a consumer of notifications delivered by an [Observable].
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

// Emit passes n to sink.
func (sink Observer[T]) Emit(n Notification[T]) {
	sink(n)
}

// ElementsOnly passes n to sink if n represents a value.
func (sink Observer[T]) ElementsOnly(n Notification[T]) {
	if n.Kind == KindNext {
		sink(n)
	}
}

// OnLastNotification creates an Observer that passes incoming emissions to
// sink, and when a notification of error or completion passes in, calls f
// just before passing it to sink.
func (sink Observer[T]) OnLastNotification(f func()) Observer[T] {
	return func(n Notification[T]) {
		switch n.Kind {
		case KindError, KindComplete:
			f()
		}

		sink(n)
	}
}

// Serialized creates an Observer that passes incoming emissions to sink
// in a mutually exclusive way.
func (sink Observer[T]) Serialized() Observer[T] {
	c := make(chan Observer[T], 1)
	c <- sink

	return func(n Notification[T]) {
		if sink, ok := <-c; ok {
			switch n.Kind {
			case KindNext:
				sink(n)
				c <- sink
			case KindError, KindComplete:
				close(c)
				sink(n)
			default: // Unknown kind.
				c <- sink
			}
		}
	}
}

// Noop gives you an Observer that does nothing.
func Noop[T any](Notification[T]) {}

// NewObserver creates an Observer from f.
func NewObserver[T any](f func(n Notification[T])) Observer[T] { return f }
