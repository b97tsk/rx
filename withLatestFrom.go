package rx

import (
	"context"
)

type withLatestFromOperator struct {
	observables []Observable
}

type withLatestFromValue struct {
	Index int
	Notification
}

func (op withLatestFromOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()

	length := len(op.observables)
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

					sink.Next(append([]interface{}(nil), values...))

				case t.HasError:
					sink(t.Notification)
					cancel()
					return

				default:
					if index > 0 {
						break
					}

					sink(t.Notification)
					cancel()
					return
				}
			}
		}
	}()

	for index, obsv := range op.observables {
		index := index
		obsv.Subscribe(ctx, func(t Notification) {
			select {
			case <-done:
			case q <- withLatestFromValue{index, t}:
			}
		})
	}

	return ctx, cancel
}

// WithLatestFrom combines the source Observable with other Observables to
// create an Observable that emits the latest values of each as a slice, only
// when the source emits.
//
// To ensure output slice has always the same length, WithLatestFrom will
// actually wait for all input Observables to emit at least once, before it
// starts emitting results.
func (o Observable) WithLatestFrom(observables ...Observable) Observable {
	observables = append([]Observable{o}, observables...)
	op := withLatestFromOperator{observables}
	return Observable{}.Lift(op.Call)
}
