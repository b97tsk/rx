package rx

import (
	"context"
)

type bufferCountOperator struct {
	bufferSize       int
	startBufferEvery int
}

func (op bufferCountOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	var (
		buffer    = make([]interface{}, 0, op.bufferSize)
		skipCount int
	)
	return source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			if skipCount > 0 {
				skipCount--
				break
			}
			buffer = append(buffer, t.Value)
			if len(buffer) < op.bufferSize {
				break
			}
			newBuffer := make([]interface{}, 0, op.bufferSize)
			if op.startBufferEvery < op.bufferSize {
				newBuffer = append(newBuffer, buffer[op.startBufferEvery:]...)
			} else {
				skipCount = op.startBufferEvery - op.bufferSize
			}
			ob.Next(buffer)
			buffer = newBuffer
		case t.HasError:
			t.Observe(ob)
		default:
			if len(buffer) > 0 {
				for op.startBufferEvery < len(buffer) {
					newBuffer := append([]interface{}(nil), buffer[op.startBufferEvery:]...)
					ob.Next(buffer)
					buffer = newBuffer
				}
				ob.Next(buffer)
			}
			t.Observe(ob)
		}
	})
}

// BufferCount buffers the source Observable values until the size hits the
// maximum bufferSize given.
//
// BufferCount collects values from the past as a slice, and emits that slice
// only when its size reaches bufferSize.
func (o Observable) BufferCount(bufferSize, startBufferEvery int) Observable {
	op := bufferCountOperator{bufferSize, startBufferEvery}
	return o.Lift(op.Call)
}
