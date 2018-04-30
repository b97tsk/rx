package rx

import (
	"context"
	"sync"
)

type exhaustMapOperator struct {
	source  Operator
	project func(interface{}, int) Observable
}

func (op exhaustMapOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()
	mu := sync.Mutex{}
	outerIndex := -1
	isActive := false
	completeSignal := make(chan struct{}, 1)
	ob = Normalize(ob)

	// Go statement makes this operator non-blocking.
	go op.source.Call(ctx, ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			mu.Lock()
			defer mu.Unlock()

			if isActive {
				break
			}

			isActive = true

			outerValue := t.Value
			outerIndex++
			outerIndex := outerIndex

			obsv := op.project(outerValue, outerIndex)

			go obsv.Subscribe(ctx, ObserverFunc(func(t Notification) {
				switch {
				case t.HasValue:
					ob.Next(t.Value)

				case t.HasError:
					ob.Error(t.Value.(error))
					cancel()

				default:
					mu.Lock()
					isActive = false
					mu.Unlock()
					select {
					case completeSignal <- struct{}{}:
					default:
					}
				}
			}))

		case t.HasError:
			ob.Error(t.Value.(error))
			cancel()

		default:
			mu.Lock()
			for isActive {
				mu.Unlock()
				select {
				case <-done:
					return
				case <-completeSignal:
				}
				mu.Lock()
			}
			mu.Unlock()
			ob.Complete()
			cancel()
		}
	}))

	return ctx, cancel
}

// Exhaust converts a higher-order Observable into a first-order Observable
// by dropping inner Observables while the previous inner Observable has not
// yet completed.
//
// Exhaust flattens an Observable-of-Observables by dropping the next inner
// Observables while the current inner is still executing.
func (o Observable) Exhaust() Observable {
	op := exhaustMapOperator{
		source:  o.Op,
		project: projectToObservable,
	}
	return Observable{op}
}

// ExhaustMap creates an Observable that projects each source value to an
// Observable which is merged in the output Observable only if the previous
// projected Observable has completed.
//
// ExhaustMap maps each value to an Observable, then flattens all of these
// inner Observables using Exhaust.
func (o Observable) ExhaustMap(project func(interface{}, int) Observable) Observable {
	op := exhaustMapOperator{
		source:  o.Op,
		project: project,
	}
	return Observable{op}
}
