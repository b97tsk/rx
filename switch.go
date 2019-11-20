package rx

import (
	"context"
)

type switchMapOperator struct {
	Project func(interface{}, int) Observable
}

func (op switchMapOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	_, childCancel := Done()

	sink = Mutex(Finally(sink, cancel))

	type X struct {
		Index           int
		ActiveIndex     int
		SourceCompleted bool
	}
	cx := make(chan *X, 1)
	cx <- &X{ActiveIndex: -1}

	source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			x := <-cx

			outerIndex := x.Index
			outerValue := t.Value
			x.Index++

			x.ActiveIndex = outerIndex

			cx <- x

			childCancel()

			obs := op.Project(outerValue, outerIndex)

			_, childCancel = obs.Subscribe(ctx, func(t Notification) {
				switch {
				case t.HasValue || t.HasError:
					sink(t)
				default:
					x := <-cx
					if x.ActiveIndex == outerIndex {
						x.ActiveIndex = -1
						if x.SourceCompleted {
							sink(t)
						}
					}
					cx <- x
				}
			})

		case t.HasError:
			sink(t)

		default:
			x := <-cx
			x.SourceCompleted = true
			if x.ActiveIndex == -1 {
				sink(t)
			}
			cx <- x
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
