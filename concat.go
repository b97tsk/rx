package rx

// Concat creates an output Observable which sequentially emits all values
// from given Observable and then moves on to the next.
//
// Concat concatenates multiple Observables together by sequentially emitting
// their values, one Observable after the other.
func Concat(observables ...Observable) Observable {
	return FromObservables(observables).Pipe(operators.ConcatAll())
}

// ConcatAll converts a higher-order Observable into a first-order Observable
// by concatenating the inner Observables in order.
//
// ConcatAll flattens an Observable-of-Observables by putting one inner
// Observable after the other.
func (Operators) ConcatAll() OperatorFunc {
	return operators.ConcatMap(ProjectToObservable)
}

// ConcatMap projects each source value to an Observable which is merged in
// the output Observable, in a serialized fashion waiting for each one to
// complete before merging the next.
//
// ConcatMap maps each value to an Observable, then flattens all of these inner
// Observables using ConcatAll.
func (Operators) ConcatMap(project func(interface{}, int) Observable) OperatorFunc {
	return func(source Observable) Observable {
		op := MergeOperator{project, 1}
		return source.Lift(op.Call)
	}
}

// ConcatMapTo projects each source value to the same Observable which is
// merged multiple times in a serialized fashion on the output Observable.
//
// It's like ConcatMap, but maps each value always to the same inner Observable.
func (Operators) ConcatMapTo(inner Observable) OperatorFunc {
	return operators.ConcatMap(func(interface{}, int) Observable { return inner })
}
