package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/queue"
)

type congestObservable struct {
	Source     rx.Observable
	BufferSize int
}

func (obs congestObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	done := ctx.Done()

	c := make(chan rx.Notification)
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

	q := make(chan rx.Notification)
	go func() {
		var queue queue.Queue
		for {
			var (
				in       <-chan rx.Notification
				out      chan<- rx.Notification
				outValue rx.Notification
			)
			length := queue.Len()
			if length < obs.BufferSize {
				in = q
			}
			if length > 0 {
				out = c
				outValue = queue.Front().(rx.Notification)
			}
			select {
			case <-done:
				return
			case t := <-in:
				queue.Push(t)
			case out <- outValue:
				queue.Pop()
			}
		}
	}()

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		select {
		case <-done:
		case q <- t:
		}
	})
}

// Congest creates an Observable that mirrors the source Observable, caches
// emissions if the source emits too fast, and congests the source if the cache
// is full.
func Congest(bufferSize int) rx.Operator {
	if bufferSize < 1 {
		return noop
	}
	return func(source rx.Observable) rx.Observable {
		return congestObservable{source, bufferSize}.Subscribe
	}
}
