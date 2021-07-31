package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/queue"
)

// Congest mirrors the source, buffers emissions if the source emits too fast,
// and blocks the source if the buffer is full.
//
// Congest has no effect if bufferSize < 1.
func Congest(bufferSize int) rx.Operator {
	if bufferSize < 1 {
		return noop
	}

	return func(source rx.Observable) rx.Observable {
		return congestObservable{source, bufferSize}.Subscribe
	}
}

type congestObservable struct {
	Source     rx.Observable
	BufferSize int
}

func (obs congestObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()

	sink = sink.WithCancel(cancel)

	cout := make(chan rx.Notification)

	go func() {
		for {
			select {
			case <-done:
				return

			case t := <-cout:
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

	cin := make(chan rx.Notification)

	go func() {
		var queue queue.Queue

		for {
			var (
				in   <-chan rx.Notification
				out  chan<- rx.Notification
				outv rx.Notification
			)

			length := queue.Len()

			if length < obs.BufferSize {
				in = cin
			}

			if length > 0 {
				out, outv = cout, queue.Front().(rx.Notification)
			}

			select {
			case <-done:
				return
			case t := <-in:
				queue.Push(t)
			case out <- outv:
				queue.Pop()
			}
		}
	}()

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		select {
		case <-done:
		case cin <- t:
		}
	})
}
