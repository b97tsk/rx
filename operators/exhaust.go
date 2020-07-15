package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/atomic"
)

type exhaustMapObservable struct {
	Source  rx.Observable
	Project func(interface{}, int) rx.Observable
}

func (obs exhaustMapObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)
	sink = sink.WithCancel(cancel).Mutex()

	sourceIndex := -1
	active := atomic.Uint32(1)

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		switch {
		case t.HasValue:
			sourceIndex++

			if !active.Cas(1, 2) {
				break
			}

			obs1 := obs.Project(t.Value, sourceIndex)

			obs1.Subscribe(ctx, func(t rx.Notification) {
				if t.HasValue || t.HasError {
					sink(t)
					return
				}
				if active.Sub(1) == 0 {
					sink(t)
				}
			})

		case t.HasError:
			sink(t)

		default:
			if active.Sub(1) == 0 {
				sink(t)
			}
		}
	})
}

// ExhaustAll converts a higher-order Observable into a first-order Observable
// by dropping inner Observables while the previous inner Observable has not
// yet completed.
//
// ExhaustAll flattens an Observable-of-Observables by dropping the next inner
// Observables while the current inner is still executing.
func ExhaustAll() rx.Operator {
	return ExhaustMap(projectToObservable)
}

// ExhaustMap creates an Observable that projects each source value to an
// Observable which is merged in the output Observable only if the previous
// projected Observable has completed.
//
// ExhaustMap maps each value to an Observable, then flattens all of these
// inner Observables using Exhaust.
func ExhaustMap(project func(interface{}, int) rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return exhaustMapObservable{source, project}.Subscribe
	}
}
