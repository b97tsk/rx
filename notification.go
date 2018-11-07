package rx

// Notification is the representation of an emission.
type Notification struct {
	Value    interface{}
	HasValue bool
	HasError bool
}

// Observe passes this Notification to the specified Observer.
func (t Notification) Observe(sink Observer) {
	sink(t)
}
