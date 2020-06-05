package rx

// GroupedObservable is an Observable type used by GroupBy operator.
type GroupedObservable struct {
	Observable
	Key interface{}
}
