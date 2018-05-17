package rx

import (
	"context"
)

type forkJoinOperator struct {
	observables []Observable
}

type forkJoinValue struct {
	Index int
	Notification
}

func (op forkJoinOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()

	length := len(op.observables)
	q := make(chan forkJoinValue, length)

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

					if !hasValues[index] {
						hasValues[index] = true
						hasValuesCount++
					}

				case t.HasError:
					ob.Error(t.Value.(error))
					cancel()
					return

				default:
					completeCount++

					if !hasValues[index] {
						ob.Complete()
						cancel()
						return
					}

					if completeCount < length {
						break
					}

					ob.Next(values)
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
			case q <- forkJoinValue{index, t}:
			}
		})
	}

	return ctx, cancel
}

// ForkJoin creates an Observable that joins last values emitted by passed
// Observables.
//
// ForkJoin waits for Observables to complete and then combine last values
// they emitted.
func ForkJoin(observables ...Observable) Observable {
	if len(observables) == 0 {
		return Empty()
	}
	op := forkJoinOperator{observables}
	return Observable{}.Lift(op.Call)
}
