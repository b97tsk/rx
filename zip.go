package rx

import (
	"container/list"
	"context"
)

type zipOperator struct {
	Observables []Observable
}

type zipValue struct {
	Index int
	Notification
}

func (op zipOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()

	sink = Finally(sink, cancel)

	length := len(op.Observables)
	q := make(chan zipValue, length)

	go func() {
		values := make([]list.List, length)
		hasValues := make([]bool, length)
		hasValuesCount := 0
		hasCompleted := make([]bool, length)
		for {
			select {
			case <-done:
				return
			case t := <-q:
				index := t.Index
				switch {
				case t.HasValue:
					values[index].PushBack(t.Value)

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

					nextValues := make([]interface{}, length)
					shouldComplete := false

					for i := range values {
						ls := &values[i]
						nextValues[i] = ls.Remove(ls.Front())
						if ls.Len() == 0 {
							hasValues[i] = false
							hasValuesCount--
							if hasCompleted[i] {
								shouldComplete = true
							}
						}
					}

					sink.Next(nextValues)

					if shouldComplete {
						sink.Complete()
						return
					}

				case t.HasError:
					sink(t.Notification)
					return

				default:
					hasCompleted[index] = true
					if !hasValues[index] {
						sink(t.Notification)
						return
					}
				}
			}
		}
	}()

	for index, obsv := range op.Observables {
		index := index
		obsv.Subscribe(ctx, func(t Notification) {
			select {
			case <-done:
			case q <- zipValue{index, t}:
			}
		})
	}

	return ctx, cancel
}

type zipAllOperator struct{}

func (op zipAllOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
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

			zip := zipOperator{observables}
			zip.Call(ctx, sink, Observable{})

		case t.HasError:
			sink(t)

		default:
			// do nothing
		}
	}, source)

	return ctx, cancel
}

// Zip combines multiple Observables to create an Observable that emits the
// values of each of its input Observables as a slice.
func Zip(observables ...Observable) Observable {
	if len(observables) == 0 {
		return Empty()
	}
	op := zipOperator{observables}
	return Observable{}.Lift(op.Call)
}

// ZipAll converts a higher-order Observable into a first-order Observable by
// waiting for the outer Observable to complete, then applying Zip.
//
// ZipAll flattens an Observable-of-Observables by applying Zip when the
// Observable-of-Observables completes.
func (Operators) ZipAll() OperatorFunc {
	return func(source Observable) Observable {
		op := zipAllOperator{}
		return source.Lift(op.Call)
	}
}
