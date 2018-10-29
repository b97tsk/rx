package rx

// StartWith creates an Observable that emits the items you specify as
// arguments before it begins to emit items emitted by the source Observable.
func (Operators) StartWith(values ...interface{}) OperatorFunc {
	return func(source Observable) Observable {
		return Concat(FromSlice(values), source)
	}
}
