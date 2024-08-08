package rx

// BufferCount buffers a number of values from the source [Observable] as
// a slice, and emits that slice when its size reaches given BufferSize;
// then, BufferCount starts a new buffer by dropping a number of most dated
// values specified by StartBufferEvery option (defaults to BufferSize).
//
// For reducing allocations, slices emitted by the output [Observable] share
// a same underlying array.
func BufferCount[T any](bufferSize int) BufferCountOperator[T] {
	return BufferCountOperator[T]{
		ts: bufferCountConfig{
			bufferSize:       bufferSize,
			startBufferEvery: bufferSize,
		},
	}
}

type bufferCountConfig struct {
	bufferSize       int
	startBufferEvery int
}

// BufferCountOperator is an [Operator] type for [BufferCount].
type BufferCountOperator[T any] struct {
	ts bufferCountConfig
}

// WithStartBufferEvery sets StartBufferEvery option to a given value.
func (op BufferCountOperator[T]) WithStartBufferEvery(n int) BufferCountOperator[T] {
	op.ts.startBufferEvery = n
	return op
}

// Apply implements the [Operator] interface.
func (op BufferCountOperator[T]) Apply(source Observable[T]) Observable[[]T] {
	if op.ts.bufferSize <= 0 {
		return Oops[[]T]("BufferCount: BufferSize <= 0")
	}

	if op.ts.startBufferEvery <= 0 {
		return Oops[[]T]("BufferCount: StartBufferEvery <= 0")
	}

	return bufferCountObservable[T]{source, op.ts}.Subscribe
}

type bufferCountObservable[T any] struct {
	source Observable[T]
	bufferCountConfig
}

func (ob bufferCountObservable[T]) Subscribe(c Context, o Observer[[]T]) {
	s := make([]T, 0, ob.bufferSize)
	skip := 0

	ob.source.Subscribe(c, func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
			if skip > 0 {
				skip--
				return
			}

			s = append(s, n.Value)

			if len(s) < ob.bufferSize {
				return
			}

			o.Next(s)

			if ob.startBufferEvery < ob.bufferSize {
				s = append(s[:0], s[ob.startBufferEvery:]...)
			} else {
				s = s[:0]
				skip = ob.startBufferEvery - ob.bufferSize
			}

		case KindComplete:
			if len(s) != 0 {
				for {
					Try1(o, Next(s), func() { o.Stop(ErrOops) })

					if len(s) <= ob.startBufferEvery {
						break
					}

					s = s[ob.startBufferEvery:]
				}
			}

			o.Complete()

		case KindError:
			o.Error(n.Error)

		case KindStop:
			o.Stop(n.Error)
		}
	})
}
