package rx

import (
	"context"
)

type combineLatestObservable struct {
	Observables []Observable
}

type combineLatestValue struct {
	Index int
	Notification
}

func (obs combineLatestObservable) Subscribe(ctx context.Context, sink Observer) {
	done := ctx.Done()
	q := make(chan combineLatestValue)

	go func() {
		length := len(obs.Observables)
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

					sink.Next(values)

				case t.HasError:
					sink(t.Notification)
					return

				default:
					if hasValues[index] {
						completeCount++
						if completeCount < length {
							break
						}
					}
					sink(t.Notification)
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
			case q <- combineLatestValue{index, t}:
			}
		})
	}
}

// CombineLatest combines multiple Observables to create an Observable that
// emits the latest values of each of its input Observables as a slice.
//
// To ensure output slice has always the same length, CombineLatest will
// actually wait for all input Observables to emit at least once, before it
// starts emitting results.
//
// For the purpose of allocation avoidance, slices emitted by the output
// Observable actually share the same underlying array.
func CombineLatest(observables ...Observable) Observable {
	if len(observables) == 0 {
		return Empty()
	}
	obs := combineLatestObservable{observables}
	return Create(obs.Subscribe)
}
