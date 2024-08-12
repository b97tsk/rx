package rx

// Range creates an [Observable] that emits a sequence of integers
// within a specified range.
func Range[T Integer](low, high T) Observable[T] {
	return func(c Context, o Observer[T]) {
		done := c.Done()

		for v := low; v < high; v++ {
			select {
			default:
			case <-done:
				o.Stop(c.Cause())
				return
			}

			Try1(o, Next(v), func() { o.Stop(ErrOops) })
		}

		o.Complete()
	}
}

// Iota creates an [Observable] that emits an infinite sequence of integers
// starting from init.
func Iota[T Integer](init T) Observable[T] {
	return func(c Context, o Observer[T]) {
		done := c.Done()

		for v := init; ; v++ {
			select {
			default:
			case <-done:
				o.Stop(c.Cause())
				return
			}

			Try1(o, Next(v), func() { o.Stop(ErrOops) })
		}
	}
}
