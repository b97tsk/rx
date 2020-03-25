package rx

import (
	"context"
)

type forkJoinObservable struct {
	Observables []Observable
}

type forkJoinValue struct {
	Index int
	Notification
}

func (obs forkJoinObservable) Subscribe(parent context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx := NewContext(parent)
	done := ctx.Done()

	sink = DoAtLast(sink, ctx.AtLast)

	length := len(obs.Observables)
	q := make(chan forkJoinValue, length)

	go func() {
		values := make([]interface{}, length)
		hasValues := make([]bool, length)
		hasValuesCount := 0
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

					if !hasValues[index] {
						hasValues[index] = true
						hasValuesCount++
					}

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

	for index, obs := range obs.Observables {
		index := index
		go obs.Subscribe(ctx, func(t Notification) {
			select {
			case <-done:
			case q <- forkJoinValue{index, t}:
			}
		})
	}

	return ctx, ctx.Cancel
}

// ForkJoin creates an Observable that joins last values emitted by passed
// Observables.
//
// ForkJoin waits for Observables to complete and then combine last values
// they emitted.
func ForkJoin(observables ...Observable) Observable {
	if len(observables) == 0 {
		return Empty()
	}
	return forkJoinObservable{observables}.Subscribe
}
