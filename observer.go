package rx

// An Observer is a consumer of values delivered by an Observable. Observers
// are simply a set of callbacks, one for each type of notification delivered
// by the Observable: Next, Error, and Complete.
type Observer interface {
	Next(interface{})
	Error(error)
	Complete()
}

// ObserverFunc is a helper type that lets you easily create an Observer from
// a function which takes a Notification as the sole argument.
type ObserverFunc func(Notification)

// Next calls the underlying function with a Next notification as argument.
func (f ObserverFunc) Next(val interface{}) {
	f(Notification{Value: val, HasValue: true})
}

// Error calls the underlying function with a Error notification as argument.
func (f ObserverFunc) Error(err error) {
	f(Notification{Value: err, HasError: true})
}

// Complete calls the underlying function with a Complete notification as argument.
func (f ObserverFunc) Complete() {
	f(Notification{})
}

// NopObserver is an Observer that does nothing.
var NopObserver Observer = ObserverFunc(func(Notification) {})

// Normalize normalizes an Observer and returns the normalized one. A normalized
// Observer ignores subsequent notifications after it receives a Complete or
// Error notification.
//
// The returned Observer is also thread-safe.
func Normalize(ob Observer) Observer {
	try := cancellableLocker{}
	return ObserverFunc(func(t Notification) {
		if try.Lock() {
			switch {
			case t.HasValue:
				ob.Next(t.Value)
				try.Unlock()
			case t.HasError:
				try.Cancel()
				ob.Error(t.Value.(error))
			default:
				try.Cancel()
				ob.Complete()
			}
		}
	})
}

func withFinalizer(ob Observer, finalize func()) ObserverFunc {
	return ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			ob.Next(t.Value)
		case t.HasError:
			ob.Error(t.Value.(error))
			finalize()
		default:
			ob.Complete()
			finalize()
		}
	})
}
