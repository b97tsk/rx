package rx

// RepeatForever repeats the stream of values emitted by the source Observable
// forever.
//
// RepeatForever does not repeat after context cancellation.
func RepeatForever[T any]() Operator[T, T] {
	return Repeat[T](-1)
}

// Repeat repeats the stream of values emitted by the source Observable
// at most count times.
//
// Repeat(0) results in an empty Observable; Repeat(1) is a no-op.
//
// Repeat does not repeat after context cancellation.
func Repeat[T any](count int) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			if count == 0 {
				return Empty[T]()
			}

			if count == 1 {
				return source
			}

			if count > 0 {
				count--
			}

			return repeatObservable[T]{source, count}.Subscribe
		},
	)
}

type repeatObservable[T any] struct {
	Source Observable[T]
	Count  int
}

func (ob repeatObservable[T]) Subscribe(c Context, o Observer[T]) {
	var observer Observer[T]

	done := c.Done()

	subscribeToSource := resistReentrance(func() {
		select {
		default:
		case <-done:
			o.Error(c.Err())
			return
		}

		ob.Source.Subscribe(c, observer)
	})

	count := ob.Count

	observer = func(n Notification[T]) {
		if n.Kind != KindComplete || count == 0 {
			o.Emit(n)
			return
		}

		if count > 0 {
			count--
		}

		subscribeToSource()
	}

	subscribeToSource()
}
