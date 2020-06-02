package rx

// An Observer is a consumer of notifications delivered by an Observable.
type Observer func(Notification)

// Next passes a NEXT notification to this Observer.
func (sink Observer) Next(val interface{}) {
	sink(Notification{Value: val, HasValue: true})
}

// Error passes an ERROR notification to this Observer.
func (sink Observer) Error(err error) {
	sink(Notification{Error: err, HasError: true})
}

// Complete passes a COMPLETE notification to this Observer.
func (sink Observer) Complete() {
	sink(Notification{})
}

// Sink passes a notification to this Observer.
//
// Note that the receiver of this method is a pointer to an Observer, this is
// useful in some cases when you need to change the receiver from one to another
// Observer, given that sink.Sink yields an Observer that is equivalent to:
//
// 	func(t Notification) { (*sink)(t) }
//
// Another use case: when you name your Observer observer, it looks bad if you
// call it like this: observer(t). Better use observer.Sink(t) instead.
//
func (sink *Observer) Sink(t Notification) {
	(*sink)(t)
}

// Noop gives you an Observer that does nothing.
func Noop(Notification) {}
