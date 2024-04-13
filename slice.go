package rx

// FromSlice creates an Observable that emits values from a slice, one after
// the other, and then completes.
func FromSlice[S ~[]T, T any](s S) Observable[T] {
	return func(c Context, o Observer[T]) {
		done := c.Done()

		for _, v := range s {
			select {
			default:
			case <-done:
				o.Error(c.Cause())
				return
			}

			Try1(o, Next(v), func() { o.Error(ErrOops) })
		}

		o.Complete()
	}
}

// Just creates an Observable that emits some values you specify as arguments,
// one after the other, and then completes.
func Just[T any](s ...T) Observable[T] {
	return FromSlice(s)
}
