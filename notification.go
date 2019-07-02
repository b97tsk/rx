package rx

// Notification is the representation of an emission.
type Notification struct {
	Value    interface{}
	Error    error
	HasValue bool
	HasError bool
}

// Observe passes this notification to the specified Observer.
func (t Notification) Observe(sink Observer) {
	sink(t)
}
