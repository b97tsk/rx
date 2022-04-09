package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/atomic"
	"github.com/b97tsk/rx/internal/critical"
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

	var cs critical.Section

	sinkAndDone := func(t rx.Notification) {
		if critical.Enter(&cs) {
			critical.Close(&cs)
			cancel()
			sink(t)
		}
	}

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

			ctx, cancel := context.WithCancel(ctx)
			childCancel = cancel

			obs1 := obs.Project(t.Value, sourceIndex)

			obs1.Subscribe(ctx, func(t rx.Notification) {
				if t.HasValue {
					if critical.Enter(&cs) {
						defer critical.Leave(&cs)

						if activeIndex.Equals(int64(sourceIndex)) {
							sink(t)
						}
					}

					return
				}

				cancel()

				if activeIndex.Cas(int64(sourceIndex), -1) {
					if t.HasError || sourceCompleted.Equals(1) && activeIndex.Equals(-1) {
						sinkAndDone(t)
					}
				}
			})

		case t.HasError:
			sinkAndDone(t)

		default:
			sourceCompleted.Store(1)

			if activeIndex.Equals(-1) {
				sinkAndDone(t)
			}
		}
	})
}
