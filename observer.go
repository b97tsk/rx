package rx

// An Observer is a consumer of notifications delivered by an Observable.
type Observer func(Notification)

// Next passes a NEXT emission to sink.
func (sink Observer) Next(val interface{}) {
	sink(Notification{Value: val, HasValue: true})
}

// Error passes an ERROR emission to sink.
func (sink Observer) Error(err error) {
	sink(Notification{Error: err, HasError: true})
}

// Complete passes a COMPLETE emission to sink.
func (sink Observer) Complete() {
	sink(Notification{})
}

// Sink passes a specified emission to *sink.
//
// Sink also yields an Observer that is equivalent to:
//
// 	func(t Notification) { (*sink)(t) }
//
// Useful when you want to set *sink to another Observer at some point.
//
// Another use case: when you name your Observer observer, it looks bad if you
// call it like this: observer(t). Better use observer.Sink(t) instead.
//
func (sink *Observer) Sink(t Notification) {
	(*sink)(t)
}

// Noop gives you an Observer that does nothing.
func Noop(Notification) {}
