package rx

// FromSlice creates an Observable that emits values from a slice, one after
// the other, and then completes.
func FromSlice[S ~[]T, T any](s S) Observable[T] {
	return func(c Context, sink Observer[T]) {
		done := c.Done()

		for _, v := range s {
			select {
			default:
			case <-done:
				sink.Error(c.Err())
				return
			}

			Try1(sink, Next(v), func() { sink.Error(ErrOops) })
		}

		sink.Complete()
	}
}

// Just creates an Observable that emits some values you specify as arguments,
// one after the other, and then completes.
func Just[T any](s ...T) Observable[T] {
	return FromSlice(s)
}
