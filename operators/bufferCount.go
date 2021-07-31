package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// BufferCount buffers the source values until the size hits the maximum
// bufferSize given.
//
// BufferCount collects values from the past as a slice, and emits that slice
// only when its size reaches bufferSize.
//
// For the purpose of allocation avoidance, slices emitted by the output
// Observable actually share a same underlying array.
func BufferCount(bufferSize int) rx.Operator {
	return BufferCountConfig{BufferSize: bufferSize}.Make()
}

// A BufferCountConfig is a configuration for BufferCount.
type BufferCountConfig struct {
	BufferSize       int
	StartBufferEvery int
}

// Make creates an Operator from this configuration.
func (config BufferCountConfig) Make() rx.Operator {
	if config.BufferSize <= 0 {
		panic("BufferCount: BufferSize is negative or zero")
	}

	if config.StartBufferEvery < 0 {
		panic("BufferCount: StartBufferEvery is negative")
	}

	if config.StartBufferEvery == 0 {
		config.StartBufferEvery = config.BufferSize
	}

	return func(source rx.Observable) rx.Observable {
		return bufferCountObservable{source, config}.Subscribe
	}
}

type bufferCountObservable struct {
	Source rx.Observable
	BufferCountConfig
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
