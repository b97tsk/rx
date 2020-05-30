package rx

import (
	"context"

	"github.com/b97tsk/rx/x/queue"
)

type zipObservable struct {
	Observables []Observable
}

type zipValue struct {
	Index int
	Notification
}

func (obs zipObservable) Subscribe(ctx context.Context, sink Observer) {
	done := ctx.Done()
	q := make(chan zipValue)

	go func() {
		type Stream struct {
			queue.Queue
			Completed bool
		}
		length := len(obs.Observables)
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
					stream.PushBack(t.Value)

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
						values[i] = stream.PopFront()
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

	for index, obs := range obs.Observables {
		index := index
		go obs.Subscribe(ctx, func(t Notification) {
			select {
			case <-done:
			case q <- zipValue{index, t}:
			}
		})
	}
}

// Zip combines multiple Observables to create an Observable that emits the
// values of each of its input Observables as a slice.
//
// For the purpose of allocation avoidance, slices emitted by the output
// Observable actually share the same underlying array.
func Zip(observables ...Observable) Observable {
	if len(observables) == 0 {
		return Empty()
	}
	obs := zipObservable{observables}
	return Create(obs.Subscribe)
}
