package rx

// StartWith creates an Observable that emits the items you specify as
// arguments before it begins to emit items emitted by the source Observable.
func (o Observable) StartWith(values ...interface{}) Observable {
	return Concat(FromSlice(values), o)
}
