package rx

import (
	"context"
	"sync"
)

type exhaustMapOperator struct {
	Project func(interface{}, int) Observable
}

func (op exhaustMapOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()
	sink = Mutex(sink)

	var (
		mu             sync.Mutex
		outerIndex     = -1
		isActive       bool
		completeSignal = make(chan struct{}, 1)
	)

	source.Subscribe(ctx, func(t Notification) {
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

			obsv := op.Project(outerValue, outerIndex)

			go obsv.Subscribe(ctx, func(t Notification) {
				switch {
				case t.HasValue:
					sink(t)

				case t.HasError:
					sink(t)
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
			})

		case t.HasError:
			sink(t)
			cancel()

		default:
			mu.Lock()
			if isActive {
				go func() {
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
					sink.Complete()
					cancel()
				}()
				return
			}
			mu.Unlock()
			sink(t)
			cancel()
		}
	})

	return ctx, cancel
}

// Exhaust converts a higher-order Observable into a first-order Observable
// by dropping inner Observables while the previous inner Observable has not
// yet completed.
//
// Exhaust flattens an Observable-of-Observables by dropping the next inner
// Observables while the current inner is still executing.
func (Operators) Exhaust() OperatorFunc {
	return operators.ExhaustMap(ProjectToObservable)
}

// ExhaustMap creates an Observable that projects each source value to an
// Observable which is merged in the output Observable only if the previous
// projected Observable has completed.
//
// ExhaustMap maps each value to an Observable, then flattens all of these
// inner Observables using Exhaust.
func (Operators) ExhaustMap(project func(interface{}, int) Observable) OperatorFunc {
	return func(source Observable) Observable {
		op := exhaustMapOperator{project}
		return source.Lift(op.Call)
	}
}
