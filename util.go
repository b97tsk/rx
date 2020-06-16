package rx

// GroupedObservable is an Observable type used by GroupBy operator.
type GroupedObservable struct {
	Observable
	Key interface{}
}

// Pair represents a struct of two values. Pair is used by Pairwise operator.
type Pair struct {
	First, Second interface{}
}
