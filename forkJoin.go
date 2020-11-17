package rx

import (
	"context"
)

// ForkJoin creates an Observable that joins last values emitted by passed
// Observables.
//
// ForkJoin waits for Observables to complete and then combine last values
// they emitted.
func ForkJoin(observables ...Observable) Observable {
	if len(observables) == 0 {
		return Empty()
	}
	return forkJoinObservable(observables).Subscribe
}

type forkJoinObservable []Observable

type forkJoinElement struct {
	Index int
	Notification
}

func (observables forkJoinObservable) Subscribe(ctx context.Context, sink Observer) {
	ctx, cancel := context.WithCancel(ctx)
	sink = sink.WithCancel(cancel)
	done := ctx.Done()

	q := make(chan forkJoinElement)

	go func() {
		length := len(observables)
		values := make([]interface{}, length)
		hasValues := make([]bool, length)
		completeCount := 0
		for {
			select {
			case <-done:
				return
			case t := <-q:
				index := t.Index
				switch {
				case t.HasValue:
					values[index] = t.Value
					hasValues[index] = true

				case t.HasError:
					sink(t.Notification)
					return

				default:
					completeCount++

					if !hasValues[index] {
						sink(t.Notification)
						return
					}

					if completeCount < length {
						break
					}

					sink.Next(values)
					sink.Complete()
					return
				}
			}
		}
	}()

	for i, obs := range observables {
		index := i
		go obs.Subscribe(ctx, func(t Notification) {
			select {
			case <-done:
			case q <- forkJoinElement{index, t}:
			}
		})
	}
}
