package rx

import (
	"context"
	"math"
	"sync/atomic"
)

type exhaustMapOperator struct {
	Project func(interface{}, int) Observable
}

func (op exhaustMapOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Mutex(Finally(sink, cancel))

	var (
		outerIndex  = -1
		activeCount = uint32(1)
	)

	source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			if !atomic.CompareAndSwapUint32(&activeCount, 1, 2) {
				break
			}

			outerIndex++
			outerIndex := outerIndex
			outerValue := t.Value

			obs := op.Project(outerValue, outerIndex)
			obs.Subscribe(ctx, func(t Notification) {
				switch {
				case t.HasValue || t.HasError:
					sink(t)
				default:
					if atomic.AddUint32(&activeCount, math.MaxUint32) > 0 {
						break
					}
					sink(t)
				}
			})

		case t.HasError:
			sink(t)

		default:
			if atomic.AddUint32(&activeCount, math.MaxUint32) > 0 {
				break
			}
			sink(t)
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
