package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// BufferCount buffers the source Observable values until the size hits the
// maximum bufferSize given.
//
// BufferCount collects values from the past as a slice, and emits that slice
// only when its size reaches bufferSize.
//
// For the purpose of allocation avoidance, slices emitted by the output
// Observable actually share the same underlying array.
func BufferCount(bufferSize int) rx.Operator {
	return BufferCountConfigure{BufferSize: bufferSize}.Make()
}

// A BufferCountConfigure is a configure for BufferCount.
type BufferCountConfigure struct {
	BufferSize       int
	StartBufferEvery int
}

// Make creates an Operator from this configure.
func (configure BufferCountConfigure) Make() rx.Operator {
	if configure.BufferSize <= 0 {
		panic("BufferCount: BufferSize is negative or zero")
	}

	if configure.StartBufferEvery < 0 {
		panic("BufferCount: StartBufferEvery is negative")
	}

	if configure.StartBufferEvery == 0 {
		configure.StartBufferEvery = configure.BufferSize
	}

	return func(source rx.Observable) rx.Observable {
		return bufferCountObservable{source, configure}.Subscribe
	}
}

type bufferCountObservable struct {
	Source rx.Observable
	BufferCountConfigure
}

func (obs bufferCountObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	buffer := make([]interface{}, 0, obs.BufferSize)
	skipCount := 0

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		switch {
		case t.HasValue:
			if skipCount > 0 {
				skipCount--

				break
			}

			buffer = append(buffer, t.Value)

			if len(buffer) < obs.BufferSize {
				break
			}

			sink.Next(buffer)

			if obs.StartBufferEvery < obs.BufferSize {
				buffer = append(buffer[:0], buffer[obs.StartBufferEvery:]...)
			} else {
				buffer = buffer[:0]
				skipCount = obs.StartBufferEvery - obs.BufferSize
			}

		case t.HasError:
			sink(t)

		default:
			if len(buffer) > 0 {
				for obs.StartBufferEvery < len(buffer) {
					sink.Next(buffer)
					buffer = buffer[obs.StartBufferEvery:]
				}

				sink.Next(buffer)
			}

			sink(t)
		}
	})
}
