package rx

import (
	"context"
)

// BufferCountOperator is an operator type.
type BufferCountOperator struct {
	BufferSize       int
	StartBufferEvery int
}

// MakeFunc creates an OperatorFunc from this operator.
func (op BufferCountOperator) MakeFunc() OperatorFunc {
	return MakeFunc(op.Call)
}

// Call invokes an execution of this operator.
func (op BufferCountOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	var (
		buffer    = make([]interface{}, 0, op.BufferSize)
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
			if len(buffer) < op.BufferSize {
				break
			}
			newBuffer := make([]interface{}, 0, op.BufferSize)
			if op.StartBufferEvery < op.BufferSize {
				newBuffer = append(newBuffer, buffer[op.StartBufferEvery:]...)
			} else {
				skipCount = op.StartBufferEvery - op.BufferSize
			}
			sink.Next(buffer)
			buffer = newBuffer
		case t.HasError:
			sink(t)
		default:
			if len(buffer) > 0 {
				for op.StartBufferEvery < len(buffer) {
					newBuffer := append([]interface{}(nil), buffer[op.StartBufferEvery:]...)
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
func (Operators) BufferCount(bufferSize int) OperatorFunc {
	return func(source Observable) Observable {
		op := BufferCountOperator{bufferSize, bufferSize}
		return source.Lift(op.Call)
	}
}
