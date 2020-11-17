package rx

import (
	"context"
)

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
	return combineLatestObservable(observables).Subscribe
}

type combineLatestObservable []Observable

type combineLatestElement struct {
	Index int
	Notification
}

func (observables combineLatestObservable) Subscribe(ctx context.Context, sink Observer) {
	ctx, cancel := context.WithCancel(ctx)
	sink = sink.WithCancel(cancel)
	done := ctx.Done()
	q := make(chan combineLatestElement)

	go func() {
		length := len(observables)
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

	for i, obs := range observables {
		index := i
		go obs.Subscribe(ctx, func(t Notification) {
			select {
			case <-done:
			case q <- combineLatestElement{index, t}:
			}
		})
	}
}
