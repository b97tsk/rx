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
		length := len(obs.Observables)
		values := make([]queue.Queue, length)
		hasValues := make([]bool, length)
		hasValuesCount := 0
		hasCompleted := make([]bool, length)
		for {
			select {
			case <-done:
				return
			case t := <-q:
				index := t.Index
				switch {
				case t.HasValue:
					values[index].PushBack(t.Value)

					if hasValuesCount < length {
						if hasValues[index] {
							break
						}

						hasValues[index] = true
						hasValuesCount++

						if hasValuesCount < length {
							break
						}
					}

					nextValues := make([]interface{}, length)
					shouldComplete := false

					for i := range values {
						queue := &values[i]
						nextValues[i] = queue.PopFront()
						if queue.Len() == 0 {
							hasValues[i] = false
							hasValuesCount--
							if hasCompleted[i] {
								shouldComplete = true
							}
						}
					}

					sink.Next(nextValues)

					if shouldComplete {
						sink.Complete()
						return
					}

				case t.HasError:
					sink(t.Notification)
					return

				default:
					hasCompleted[index] = true
					if !hasValues[index] {
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
func Zip(observables ...Observable) Observable {
	if len(observables) == 0 {
		return Empty()
	}
	obs := zipObservable{observables}
	return Create(obs.Subscribe)
}
