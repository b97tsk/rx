package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/atomic"
)

// SwitchAll converts a higher-order Observable into a first-order Observable
// by subscribing to only the most recently emitted of those inner Observables.
//
// SwitchAll flattens an Observable-of-Observables by dropping the previous
// inner Observable once a new one appears.
func SwitchAll() rx.Operator {
	return SwitchMap(projectToObservable)
}

// SwitchMap projects each source value to an Observable which is merged in
// the output Observable, emitting values only from the most recently projected
// Observable.
//
// SwitchMap maps each value to an Observable, then flattens all of these inner
// Observables using SwitchAll.
func SwitchMap(project func(interface{}, int) rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return switchMapObservable{source, project}.Subscribe
	}
}

// SwitchMapTo projects each source value to the same Observable which is
// flattened multiple times with SwitchAll in the output Observable.
//
// It's like SwitchMap, but maps each value always to the same inner Observable.
func SwitchMapTo(inner rx.Observable) rx.Operator {
	return SwitchMap(func(interface{}, int) rx.Observable { return inner })
}

type switchMapObservable struct {
	Source  rx.Observable
	Project func(interface{}, int) rx.Observable
}

func (obs switchMapObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel).Mutex()

	var childCancel context.CancelFunc

	activeIndex := atomic.FromInt64(-1)
	sourceIndex := -1
	sourceCompleted := atomic.FromUint32(0)

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		switch {
		case t.HasValue:
			sourceIndex++

			sourceIndex := sourceIndex
			activeIndex.Store(int64(sourceIndex))

			if childCancel != nil {
				childCancel()
			}

			var childCtx context.Context

			childCtx, childCancel = context.WithCancel(ctx)

			obs1 := obs.Project(t.Value, sourceIndex)

			obs1.Subscribe(childCtx, func(t rx.Notification) {
				if !t.HasValue {
					childCancel()
				}

				if t.HasValue || t.HasError {
					sink(t)
					return
				}

				if activeIndex.Cas(int64(sourceIndex), -1) &&
					sourceCompleted.Equals(1) && activeIndex.Equals(-1) {
					sink(t)
				}
			})

		case t.HasError:
			sink(t)

		default:
			sourceCompleted.Store(1)

			if activeIndex.Equals(-1) {
				sink(t)
			}
		}
	})
}
