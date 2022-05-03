package rx

// Notification is the representation of an emission.
//
// There is three kinds of Notifications: values, errors and completions.
//
// An Observable can only emit N+1 Notifications: either N values and an error,
// or N values and a completion. The last Notification emitted by an Observable
// must be an error or a completion.
//
type Notification[T any] struct {
	HasValue bool
	HasError bool
	Value    T
	Error    error
}

// Next creates a Notification that represents a value.
func Next[T any](v T) Notification[T] {
	return Notification[T]{Value: v, HasValue: true}
}

// Error creates a Notification that represents an error.
func Error[T any](e error) Notification[T] {
	return Notification[T]{Error: e, HasError: true}
}

// Complete creates a Notification that represents a completion.
func Complete[T any]() Notification[T] {
	return Notification[T]{}
}
