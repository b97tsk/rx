package rx

// Notification is the representation of an emission.
type Notification struct {
	HasValue bool
	HasError bool
	Value    interface{}
	Error    error
}
