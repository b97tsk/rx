package rx

// Notification is the representation of an emission.
type Notification struct {
	Value    interface{}
	Error    error
	HasValue bool
	HasError bool
}
