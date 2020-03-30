package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// A BufferCountConfigure is a configure for BufferCount.
type BufferCountConfigure struct {
	BufferSize       int
	StartBufferEvery int
}

// Use creates an Operator from this configure.
func (configure BufferCountConfigure) Use() rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return bufferCountObservable{source, configure}.Subscribe
	}
}

type bufferCountObservable struct {
	Source rx.Observable
	BufferCountConfigure
}

func (obs bufferCountObservable) Subscribe(ctx context.Context, sink rx.Observer) (context.Context, context.CancelFunc) {
	var (
		buffer    = make([]interface{}, 0, obs.BufferSize)
		skipCount int
	)
	if obs.StartBufferEvery == 0 {
		obs.StartBufferEvery = obs.BufferSize
	}
	return obs.Source.Subscribe(ctx, func(t rx.Notification) {
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
			newBuffer := make([]interface{}, 0, obs.BufferSize)
			if obs.StartBufferEvery < obs.BufferSize {
				newBuffer = append(newBuffer, buffer[obs.StartBufferEvery:]...)
			} else {
				skipCount = obs.StartBufferEvery - obs.BufferSize
			}
			sink.Next(buffer)
			buffer = newBuffer
		case t.HasError:
			sink(t)
		default:
			if len(buffer) > 0 {
				for obs.StartBufferEvery < len(buffer) {
					remains := buffer[obs.StartBufferEvery:]
					newBuffer := make([]interface{}, len(remains))
					copy(newBuffer, remains)
					sink.Next(buffer)
					buffer = newBuffer
				}
				sink.Next(buffer)
			}
			sink(t)
		}
	})
}

// BufferCount buffers the source Observable values until the size hits the
// maximum bufferSize given.
//
// BufferCount collects values from the past as a slice, and emits that slice
// only when its size reaches bufferSize.
func BufferCount(bufferSize int) rx.Operator {
	return BufferCountConfigure{BufferSize: bufferSize}.Use()
}
