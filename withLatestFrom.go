package rx

import (
	"context"
)

type withLatestFromObservable struct {
	Observables []Observable
}

type withLatestFromValue struct {
	Index int
	Notification
}

func (obs withLatestFromObservable) Subscribe(parent context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx := NewContext(parent)
	done := ctx.Done()

	sink = DoAtLast(sink, ctx.AtLast)

	length := len(obs.Observables)
	q := make(chan withLatestFromValue, length)

	go func() {
		values := make([]interface{}, length)
		hasValues := make([]bool, length)
		hasValuesCount := 0
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

					if index > 0 {
						break
					}

					newValues := make([]interface{}, len(values))
					copy(newValues, values)
					sink.Next(newValues)

				case t.HasError:
					sink(t.Notification)
					return

				default:
					if index > 0 {
						break
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
			case q <- withLatestFromValue{index, t}:
			}
		})
	}

	return ctx, ctx.Cancel
}

// WithLatestFrom combines the source Observable with other Observables to
// create an Observable that emits the latest values of each as a slice, only
// when the source emits.
//
// To ensure output slice has always the same length, WithLatestFrom will
// actually wait for all input Observables to emit at least once, before it
// starts emitting results.
func (Operators) WithLatestFrom(observables ...Observable) Operator {
	return func(source Observable) Observable {
		observables = append([]Observable{source}, observables...)
		return withLatestFromObservable{observables}.Subscribe
	}
}
