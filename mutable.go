package rx

// A MutableObserver is an Observer that mirrors emission to the underlying
// Observer, which can be changed programmatically.
type MutableObserver struct {
	Observer Observer
}

// Next calls the Next method of the underlying Observer.
func (ob *MutableObserver) Next(val interface{}) {
	ob.Observer.Next(val)
}

// Error calls the Error method of the underlying Observer.
func (ob *MutableObserver) Error(err error) {
	ob.Observer.Error(err)
}

// Complete calls the Complete method of the underlying Observer.
func (ob *MutableObserver) Complete() {
	ob.Observer.Complete()
}
