package rx

import (
	"context"

	"github.com/b97tsk/rx/x/atomic"
)

type exhaustMapObservable struct {
	Source  Observable
	Project func(interface{}, int) Observable
}

func (obs exhaustMapObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Mutex(Finally(sink, cancel))

	var (
		sourceIndex = -1
		activeCount = atomic.Uint32(1)
	)

	obs.Source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			if !activeCount.Cas(1, 2) {
				break
			}

			sourceIndex++
			sourceIndex := sourceIndex
			sourceValue := t.Value

			obs := obs.Project(sourceValue, sourceIndex)
			obs.Subscribe(ctx, func(t Notification) {
				switch {
				case t.HasValue || t.HasError:
					sink(t)
				default:
					if activeCount.Sub(1) > 0 {
						break
					}
					sink(t)
				}
			})

		case t.HasError:
			sink(t)

		default:
			if activeCount.Sub(1) > 0 {
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
func (Operators) Exhaust() Operator {
	return operators.ExhaustMap(ProjectToObservable)
}

// ExhaustMap creates an Observable that projects each source value to an
// Observable which is merged in the output Observable only if the previous
// projected Observable has completed.
//
// ExhaustMap maps each value to an Observable, then flattens all of these
// inner Observables using Exhaust.
func (Operators) ExhaustMap(project func(interface{}, int) Observable) Operator {
	return func(source Observable) Observable {
		return exhaustMapObservable{source, project}.Subscribe
	}
}
