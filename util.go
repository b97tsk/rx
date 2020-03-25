package rx

// ProjectToObservable type-asserts each value to be an Observable and returns
// it. If type assertion fails, it returns Throw(ErrNotObservable).
func ProjectToObservable(val interface{}, idx int) Observable {
	if obs, ok := val.(Observable); ok {
		return obs
	}
	return Throw(ErrNotObservable)
}
