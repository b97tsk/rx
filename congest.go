package rx

import (
	"context"

	"github.com/b97tsk/rx/x/queue"
)

type congestOperator struct {
	BufferSize int
}

func (op congestOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()

	sink = Finally(sink, cancel)

	c := make(chan Notification)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-c:
				switch {
				case t.HasValue:
					sink(t)
				default:
					sink(t)
					return
				}
			}
		}
	}()

	q := make(chan Notification)
	go func() {
		var queue queue.Queue
		for {
			var (
				in       <-chan Notification
				out      chan<- Notification
				outValue Notification
			)
			length := queue.Len()
			if length < op.BufferSize {
				in = q
			}
			if length > 0 {
				out = c
				outValue = queue.Front().(Notification)
			}
			select {
			case <-done:
				return
			case t := <-in:
				queue.PushBack(t)
			case out <- outValue:
				queue.PopFront()
			}
		}
	}()

	source.Subscribe(ctx, func(t Notification) {
		select {
		case <-done:
		case q <- t:
		}
	})

	return ctx, cancel
}

// Congest creates an Observable that mirrors the source Observable, caches
// emissions if the source emits too fast, and congests the source if the cache
// is full.
func (Operators) Congest(bufferSize int) OperatorFunc {
	return func(source Observable) Observable {
		if bufferSize < 1 {
			return source
		}
		op := congestOperator{bufferSize}
		return source.Lift(op.Call)
	}
}
