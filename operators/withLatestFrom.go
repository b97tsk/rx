package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

type withLatestFromObservable []rx.Observable

type withLatestFromElement struct {
	Index int
	rx.Notification
}

func (observables withLatestFromObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	done := ctx.Done()
	q := make(chan withLatestFromElement)

	go func() {
		length := len(observables)
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

					sink.Next(values)

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

	for i, obs := range observables {
		index := i
		go obs.Subscribe(ctx, func(t rx.Notification) {
			select {
			case <-done:
			case q <- withLatestFromElement{index, t}:
			}
		})
	}
}

// WithLatestFrom combines the source Observable with other Observables to
// create an Observable that emits the latest values of each as a slice, only
// when the source emits.
//
// To ensure output slice has always the same length, WithLatestFrom will
// actually wait for all input Observables to emit at least once, before it
// starts emitting results.
//
// For the purpose of allocation avoidance, slices emitted by the output
// Observable actually share the same underlying array.
func WithLatestFrom(observables ...rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		observables = append([]rx.Observable{source}, observables...)
		return rx.Create(withLatestFromObservable(observables).Subscribe)
	}
}
