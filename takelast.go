package rx

import "github.com/b97tsk/rx/internal/queue"

// TakeLast emits only the last count values emitted by the source [Observable].
func TakeLast[T any](count int) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			if count <= 0 {
				return Empty[T]()
			}

			return func(c Context, o Observer[T]) {
				var q queue.Queue[T]

				source.Subscribe(c, func(n Notification[T]) {
					switch n.Kind {
					case KindNext:
						if q.Len() == count {
							q.Pop()
						}

						q.Push(n.Value)

					case KindComplete:
						done := c.Done()

						for i := range q.Len() {
							select {
							default:
							case <-done:
								o.Stop(c.Cause())
								return
							}

							Try1(o, Next(q.At(i)), func() { o.Stop(ErrOops) })
						}

						o.Emit(n)

					case KindError, KindStop:
						o.Emit(n)
					}
				})
			}
		},
	)
}
