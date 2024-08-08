package rx

// NotificationKind is the type of an emission.
type NotificationKind int8

const (
	_ NotificationKind = iota
	KindNext
	KindComplete
	KindError
	KindStop
)

// Notification is the representation of an emission.
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
type Notification[T any] struct {
	Kind  NotificationKind
	Value T
	Error error
}

// Next creates a [Notification] that represents a value.
func Next[T any](v T) Notification[T] {
	return Notification[T]{Kind: KindNext, Value: v}
}

// Complete creates a [Notification] that represents a completion.
func Complete[T any]() Notification[T] {
	return Notification[T]{Kind: KindComplete}
}

// Error creates a [Notification] that represents an error.
func Error[T any](err error) Notification[T] {
	return Notification[T]{Kind: KindError, Error: err}
}

// Stop creates a [Notification] that represents a stop.
func Stop[T any](err error) Notification[T] {
	return Notification[T]{Kind: KindStop, Error: err}
}

// And returns n or other if n represents a value or completion.
func (n Notification[T]) And(other Notification[T]) Notification[T] {
	switch n.Kind {
	case KindNext, KindComplete:
		return other
	default:
		return n
	}
}

// Or returns n or other if n does not represent a value or completion.
func (n Notification[T]) Or(other Notification[T]) Notification[T] {
	switch n.Kind {
	case KindNext, KindComplete:
		return n
	default:
		return other
	}
}
