package rx

import (
	"context"
	"sync"
)

type switchMapOperator struct {
	Project func(interface{}, int) Observable
}

func (op switchMapOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Mutex(Finally(sink, cancel))

	var (
		childCtx    = canceledCtx
		childCancel = nothingToDo

		mutex           sync.Mutex
		outerIndex      = -1
		activeIndex     = -1
		sourceCompleted bool
	)

	source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			mutex.Lock()
			defer mutex.Unlock()

			outerIndex++
			outerIndex := outerIndex
			outerValue := t.Value

			activeIndex = outerIndex
			childCancel()

			obsv := op.Project(outerValue, outerIndex)

			childCtx, childCancel = context.WithCancel(ctx)

			go obsv.Subscribe(childCtx, func(t Notification) {
				switch {
				case t.HasValue:
					sink(t)

				case t.HasError:
					sink(t)

				default:
					mutex.Lock()
					defer mutex.Unlock()

					if activeIndex != outerIndex {
						break
					}

					activeIndex = -1

					if sourceCompleted {
						sink(t)
					}
				}
			})

		case t.HasError:
			sink(t)

		default:
			mutex.Lock()
			defer mutex.Unlock()
			sourceCompleted = true
			if activeIndex == -1 {
				sink(t)
			}
		}
	})

	return ctx, cancel
}

// Switch converts a higher-order Observable into a first-order Observable by
// subscribing to only the most recently emitted of those inner Observables.
//
// Switch flattens an Observable-of-Observables by dropping the previous inner
// Observable once a new one appears.
func (Operators) Switch() OperatorFunc {
	return operators.SwitchMap(ProjectToObservable)
}

// SwitchMap creates an Observable that projects each source value to an
// Observable which is merged in the output Observable, emitting values only
// from the most recently projected Observable.
//
// SwitchMap maps each value to an Observable, then flattens all of these inner
// Observables using Switch.
func (Operators) SwitchMap(project func(interface{}, int) Observable) OperatorFunc {
	return func(source Observable) Observable {
		op := switchMapOperator{project}
		return source.Lift(op.Call)
	}
}

// SwitchMapTo creates an Observable that projects each source value to the
// same Observable which is flattened multiple times with Switch in the output
// Observable.
//
// It's like SwitchMap, but maps each value always to the same inner Observable.
func (Operators) SwitchMapTo(inner Observable) OperatorFunc {
	return operators.SwitchMap(func(interface{}, int) Observable { return inner })
}
