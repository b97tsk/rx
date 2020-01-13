package rx

// An Operator is an operation on an Observable. When called, they do not
// change the existing Observable instance. Instead, they return a new
// Observable, whose subscription logic is based on the first Observable.
type Operator func(Observable) Observable

// Pipe stitches operators together into a chain.
func Pipe(operators ...Operator) Operator {
	return func(source Observable) Observable {
		return source.Pipe(operators...)
	}
}
