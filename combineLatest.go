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

func (obs combineLatestObservable) Subscribe(parent context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx := NewContext(parent)
	done := ctx.Done()

	sink = DoAtLast(sink, ctx.AtLast)

	length := len(obs.Observables)
	q := make(chan combineLatestValue, length)

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

					newValues := make([]interface{}, len(values))
					copy(newValues, values)
					sink.Next(newValues)

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

	return ctx, ctx.Cancel
}

// CombineLatest combines multiple Observables to create an Observable that
// emits the latest values of each of its input Observables as a slice.
//
// To ensure output slice has always the same length, CombineLatest will
// actually wait for all input Observables to emit at least once, before it
// starts emitting results.
func CombineLatest(observables ...Observable) Observable {
	if len(observables) == 0 {
		return Empty()
	}
	return combineLatestObservable{observables}.Subscribe
}

// CombineAll converts a higher-order Observable into a first-order Observable
// by waiting for the outer Observable to complete, then applying CombineLatest.
//
// CombineAll flattens an Observable-of-Observables by applying CombineLatest
// when the Observable-of-Observables completes.
func (Operators) CombineAll() Operator {
	return ToObservablesConfigure{CombineLatest}.Use()
}
