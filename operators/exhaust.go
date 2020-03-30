package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/x/atomic"
)

type exhaustMapObservable struct {
	Source  rx.Observable
	Project func(interface{}, int) rx.Observable
}

func (obs exhaustMapObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	sink = rx.Mutex(sink)

	var (
		sourceIndex = -1
		activeCount = atomic.Uint32(1)
	)

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		switch {
		case t.HasValue:
			if !activeCount.Cas(1, 2) {
				break
			}

			sourceIndex++
			sourceIndex := sourceIndex
			sourceValue := t.Value

			obs := obs.Project(sourceValue, sourceIndex)
			obs.Subscribe(ctx, func(t rx.Notification) {
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
}

// Exhaust converts a higher-order Observable into a first-order Observable
// by dropping inner Observables while the previous inner Observable has not
// yet completed.
//
// Exhaust flattens an Observable-of-Observables by dropping the next inner
// Observables while the current inner is still executing.
func Exhaust() rx.Operator {
	return ExhaustMap(rx.ProjectToObservable)
}

// ExhaustMap creates an Observable that projects each source value to an
// Observable which is merged in the output Observable only if the previous
// projected Observable has completed.
//
// ExhaustMap maps each value to an Observable, then flattens all of these
// inner Observables using Exhaust.
func ExhaustMap(project func(interface{}, int) rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := exhaustMapObservable{source, project}
		return rx.Create(obs.Subscribe)
	}
}
