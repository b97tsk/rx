package rx

import (
	"context"
)

type combineLatestOperator struct {
	Observables []Observable
}

type combineLatestValue struct {
	Index int
	Notification
}

func (op combineLatestOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()

	sink = Finally(sink, cancel)

	length := len(op.Observables)
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

					sink.Next(append([]interface{}(nil), values...))

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

	for index, obsv := range op.Observables {
		index := index
		obsv.Subscribe(ctx, func(t Notification) {
			select {
			case <-done:
			case q <- combineLatestValue{index, t}:
			}
		})
	}

	return ctx, cancel
}

type combineAllOperator struct{}

func (op combineAllOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	toObservablesOperator(op).Call(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			observables := t.Value.([]Observable)

			if len(observables) == 0 {
				sink.Complete()
				break
			}

			combineLatest := combineLatestOperator{observables}
			combineLatest.Call(ctx, sink, Observable{})

		case t.HasError:
			sink(t)

		default:
			// do nothing
		}
	}, source)

	return ctx, cancel
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
	op := combineLatestOperator{observables}
	return Observable{}.Lift(op.Call)
}

// CombineAll converts a higher-order Observable into a first-order Observable
// by waiting for the outer Observable to complete, then applying CombineLatest.
//
// CombineAll flattens an Observable-of-Observables by applying CombineLatest
// when the Observable-of-Observables completes.
func (Operators) CombineAll() OperatorFunc {
	return func(source Observable) Observable {
		op := combineAllOperator{}
		return source.Lift(op.Call)
	}
}
