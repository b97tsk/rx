package rx

import "golang.org/x/exp/constraints"

// Range creates an Observable that emits a sequence of integers
// within a specified range.
func Range[T constraints.Integer](low, high T) Observable[T] {
	return func(c Context, sink Observer[T]) {
		done := c.Done()

		for v := low; v < high; v++ {
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

// Iota creates an Observable that emits an infinite sequence of integers
// starting from init.
func Iota[T constraints.Integer](init T) Observable[T] {
	return func(c Context, sink Observer[T]) {
		done := c.Done()

		for v := init; ; v++ {
			select {
			default:
			case <-done:
				sink.Error(c.Err())
				return
			}

			Try1(sink, Next(v), func() { sink.Error(ErrOops) })
		}
	}
}
