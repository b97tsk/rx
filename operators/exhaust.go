package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/atomic"
)

// ExhaustAll converts a higher-order Observable into a first-order Observable
// by dropping inner Observables while the previous inner Observable has not
// yet completed.
//
// ExhaustAll flattens an Observable-of-Observables by dropping the next inner
// Observables while the current inner is still executing.
func ExhaustAll() rx.Operator {
	return ExhaustMap(projectToObservable)
}

// ExhaustMap projects each source value to an Observable which is merged
// in the output Observable only if the previous projected Observable has
// completed.
//
// ExhaustMap maps each value to an Observable, then flattens all of these
// inner Observables using ExhaustAll.
func ExhaustMap(project func(interface{}, int) rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return exhaustMapObservable{source, project}.Subscribe
	}
}

// ExhaustMapTo projects each source value to the same Observable which is
// flattened multiple times with ExhaustAll in the output Observable.
//
// It's like ExhaustMap, but maps each value always to the same inner
// Observable.
func ExhaustMapTo(inner rx.Observable) rx.Operator {
	return ExhaustMap(func(interface{}, int) rx.Observable { return inner })
}

type exhaustMapObservable struct {
	Source  rx.Observable
	Project func(interface{}, int) rx.Observable
}

func (obs exhaustMapObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel).Mutex()

	sourceIndex := -1
	workers := atomic.FromUint32(1)

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		switch {
		case t.HasValue:
			sourceIndex++

			if workers.Cas(1, 2) {
				obs1 := obs.Project(t.Value, sourceIndex)

				obs1.Subscribe(ctx, func(t rx.Notification) {
					if t.HasValue || t.HasError {
						sink(t)

						return
					}

					if workers.Sub(1) == 0 {
						sink(t)
					}
				})
			}

		case t.HasError:
			sink(t)

		default:
			if workers.Sub(1) == 0 {
				sink(t)
			}
		}
	})
}
