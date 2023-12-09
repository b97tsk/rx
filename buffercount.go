package rx

import (
	"context"
)

// BufferCount buffers a number of values from the source Observable as
// a slice, and emits that slice when its size reaches given BufferSize;
// then, BufferCount starts a new buffer by dropping a number of most dated
// values specified by StartBufferEvery option (defaults to BufferSize).
//
// For reducing allocations, slices emitted by the output Observable share
// a same underlying array.
func BufferCount[T any](bufferSize int) BufferCountOperator[T] {
	if bufferSize <= 0 {
		panic("bufferSize <= 0")
	}

	return BufferCountOperator[T]{
		opts: bufferCountConfig{
			BufferSize:       bufferSize,
			StartBufferEvery: bufferSize,
		},
	}
}

type bufferCountConfig struct {
	BufferSize       int
	StartBufferEvery int
}

// BufferCountOperator is an [Operator] type for [BufferCount].
type BufferCountOperator[T any] struct {
	opts bufferCountConfig
}

// WithStartBufferEvery sets StartBufferEvery option to a given value.
func (op BufferCountOperator[T]) WithStartBufferEvery(n int) BufferCountOperator[T] {
	if n <= 0 {
		panic("n <= 0")
	}

	op.opts.StartBufferEvery = n

	return op
}

// Apply implements the Operator interface.
func (op BufferCountOperator[T]) Apply(source Observable[T]) Observable[[]T] {
	return bufferCountObservable[T]{source, op.opts}.Subscribe
}

type bufferCountObservable[T any] struct {
	Source Observable[T]
	bufferCountConfig
}

func (obs bufferCountObservable[T]) Subscribe(ctx context.Context, sink Observer[[]T]) {
	s := make([]T, 0, obs.BufferSize)
	skip := 0

	obs.Source.Subscribe(ctx, func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
			if skip > 0 {
				skip--
				break
			}

			s = append(s, n.Value)

			if len(s) < obs.BufferSize {
				break
			}

			sink.Next(s)

			if obs.StartBufferEvery < obs.BufferSize {
				s = append(s[:0], s[obs.StartBufferEvery:]...)
			} else {
				s = s[:0]
				skip = obs.StartBufferEvery - obs.BufferSize
			}

		case KindError:
			sink.Error(n.Error)

		case KindComplete:
			if len(s) > 0 {
				for {
					sink.Next(s)

					if len(s) <= obs.StartBufferEvery {
						break
					}

					s = s[obs.StartBufferEvery:]
				}
			}

			sink.Complete()
		}
	})
}
