package rx

import (
	"context"

	"github.com/b97tsk/rx/internal/queue"
)

// Zip combines multiple Observables to create an Observable that emits the
// values of each of its input Observables as a slice.
//
// For the purpose of allocation avoidance, slices emitted by the output
// Observable actually share the same underlying array.
func Zip(observables ...Observable) Observable {
	if len(observables) == 0 {
		return Empty()
	}
	return zipObservable(observables).Subscribe
}

type zipObservable []Observable

type zipElement struct {
	Index int
	Notification
}

func (observables zipObservable) Subscribe(ctx context.Context, sink Observer) {
	ctx, cancel := context.WithCancel(ctx)
	sink = sink.WithCancel(cancel)
	done := ctx.Done()
	q := make(chan zipElement)

	go func() {
		type Stream struct {
			queue.Queue
			Completed bool
		}
		length := len(observables)
		streams := make([]Stream, length)
		values := make([]interface{}, length)
		readyCount := 0
		for {
			select {
			case <-done:
				return
			case t := <-q:
				index := t.Index
				switch {
				case t.HasValue:
					stream := &streams[index]
					stream.Push(t.Value)

					if readyCount < length {
						if stream.Len() > 1 {
							break
						}
						readyCount++
						if readyCount < length {
							break
						}
					}

					shouldComplete := false
					for i := range streams {
						stream := &streams[i]
						values[i] = stream.Pop()
						if stream.Len() == 0 {
							readyCount--
							if stream.Completed {
								shouldComplete = true
							}
						}
					}
					sink.Next(values)
					if shouldComplete {
						sink.Complete()
						return
					}

				case t.HasError:
					sink(t.Notification)
					return

				default:
					stream := &streams[index]
					stream.Completed = true
					if stream.Len() == 0 {
						sink(t.Notification)
						return
					}
				}
			}
		}
	}()

	for i, obs := range observables {
		index := i
		go obs.Subscribe(ctx, func(t Notification) {
			select {
			case <-done:
			case q <- zipElement{index, t}:
			}
		})
	}
}
