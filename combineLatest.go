package rx

import (
	"context"
)

type combineLatestOperator struct {
	observables []Observable
}

type combineLatestValue struct {
	Index int
	Notification
}

func (op combineLatestOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()

	length := len(op.observables)
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

					ob.Next(append([]interface{}(nil), values...))

				case t.HasError:
					ob.Error(t.Value.(error))
					cancel()
					return

				default:
					if hasValues[index] {
						completeCount++
						if completeCount < length {
							break
						}
					}
					ob.Complete()
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
			case q <- combineLatestValue{index, t}:
			}
		})
	}

	return ctx, cancel
}

type combineAllOperator struct{}

func (op combineAllOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	toObservablesOperator(op).Call(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			observables := t.Value.([]Observable)

			if len(observables) == 0 {
				ob.Complete()
				cancel()
				break
			}

			combineLatest := combineLatestOperator{observables}
			combineLatest.Call(ctx, withFinalizer(ob, cancel), Observable{})

		case t.HasError:
			ob.Error(t.Value.(error))
			cancel()

		default:
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
func (o Observable) CombineAll() Observable {
	op := combineAllOperator{}
	return o.Lift(op.Call)
}
