package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

type withLatestFromObservable struct {
	Observables []rx.Observable
}

type withLatestFromValue struct {
	Index int
	rx.Notification
}

func (obs withLatestFromObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	done := ctx.Done()

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
		go obs.Subscribe(ctx, func(t rx.Notification) {
			select {
			case <-done:
			case q <- withLatestFromValue{index, t}:
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
func WithLatestFrom(observables ...rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		observables = append([]rx.Observable{source}, observables...)
		obs := withLatestFromObservable{observables}
		return rx.Create(obs.Subscribe)
	}
}
