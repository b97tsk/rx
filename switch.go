package rx

import (
	"context"
)

type switchMapObservable struct {
	Source  Observable
	Project func(interface{}, int) Observable
}

func (obs switchMapObservable) Subscribe(ctx context.Context, sink Observer) {
	sink = Mutex(sink)

	type X struct {
		Index           int
		ActiveIndex     int
		SourceCompleted bool
	}
	cx := make(chan *X, 1)
	cx <- &X{ActiveIndex: -1}

	var childCancel context.CancelFunc

	obs.Source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			x := <-cx

			sourceIndex := x.Index
			sourceValue := t.Value
			x.Index++

			x.ActiveIndex = sourceIndex

			cx <- x

			if childCancel != nil {
				childCancel()
			}

			obs := obs.Project(sourceValue, sourceIndex)

			_, childCancel = obs.Subscribe(ctx, func(t Notification) {
				switch {
				case t.HasValue || t.HasError:
					sink(t)
				default:
					x := <-cx
					if x.ActiveIndex == sourceIndex {
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
}

// Switch converts a higher-order Observable into a first-order Observable by
// subscribing to only the most recently emitted of those inner Observables.
//
// Switch flattens an Observable-of-Observables by dropping the previous inner
// Observable once a new one appears.
func (Operators) Switch() Operator {
	return operators.SwitchMap(ProjectToObservable)
}

// SwitchMap creates an Observable that projects each source value to an
// Observable which is merged in the output Observable, emitting values only
// from the most recently projected Observable.
//
// SwitchMap maps each value to an Observable, then flattens all of these inner
// Observables using Switch.
func (Operators) SwitchMap(project func(interface{}, int) Observable) Operator {
	return func(source Observable) Observable {
		obs := switchMapObservable{source, project}
		return Create(obs.Subscribe)
	}
}

// SwitchMapTo creates an Observable that projects each source value to the
// same Observable which is flattened multiple times with Switch in the output
// Observable.
//
// It's like SwitchMap, but maps each value always to the same inner Observable.
func (Operators) SwitchMapTo(inner Observable) Operator {
	return operators.SwitchMap(func(interface{}, int) Observable { return inner })
}
