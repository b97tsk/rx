package rx

// Notification is the representation of an emission.
type Notification struct {
	Value    interface{}
	HasValue bool
	HasError bool
}

// Observe notifies this Notification to the specified Observer.
func (t Notification) Observe(ob Observer) {
	switch {
	case t.HasValue:
		ob.Next(t.Value)
	case t.HasError:
		ob.Error(t.Value.(error))
	default:
		ob.Complete()
	}
}
