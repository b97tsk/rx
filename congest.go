package rx

import (
	"context"

	"github.com/b97tsk/rx/x/queue"
)

type congestObservable struct {
	Source     Observable
	BufferSize int
}

func (obs congestObservable) Subscribe(parent context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx := NewContext(parent)
	done := ctx.Done()

	sink = DoAtLast(sink, ctx.AtLast)

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
			if length < obs.BufferSize {
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

	obs.Source.Subscribe(ctx, func(t Notification) {
		select {
		case <-done:
		case q <- t:
		}
	})

	return ctx, ctx.Cancel
}

// Congest creates an Observable that mirrors the source Observable, caches
// emissions if the source emits too fast, and congests the source if the cache
// is full.
func (Operators) Congest(bufferSize int) Operator {
	return func(source Observable) Observable {
		if bufferSize < 1 {
			return source
		}
		return congestObservable{source, bufferSize}.Subscribe
	}
}
