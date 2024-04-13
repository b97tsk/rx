package rx

import "golang.org/x/exp/constraints"

// Range creates an Observable that emits a sequence of integers
// within a specified range.
func Range[T constraints.Integer](low, high T) Observable[T] {
	return func(c Context, o Observer[T]) {
		done := c.Done()

		for v := low; v < high; v++ {
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

// Iota creates an Observable that emits an infinite sequence of integers
// starting from init.
func Iota[T constraints.Integer](init T) Observable[T] {
	return func(c Context, o Observer[T]) {
		done := c.Done()

		for v := init; ; v++ {
			select {
			default:
			case <-done:
				o.Error(c.Cause())
				return
			}

			Try1(o, Next(v), func() { o.Error(ErrOops) })
		}
	}
}
