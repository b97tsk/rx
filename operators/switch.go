package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

type switchMapObservable struct {
	Source  rx.Observable
	Project func(interface{}, int) (rx.Observable, error)
}

func (obs switchMapObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	sink = rx.Mutex(sink)

	type X struct {
		Index           int
		ActiveIndex     int
		SourceCompleted bool
	}
	cx := make(chan *X, 1)
	cx <- &X{ActiveIndex: -1}

	var childCancel context.CancelFunc

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
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

			obs, err := obs.Project(sourceValue, sourceIndex)
			if err != nil {
				sink.Error(err)
				return
			}

			_, childCancel = obs.Subscribe(ctx, func(t rx.Notification) {
				if t.HasValue || t.HasError {
					sink(t)
					return
				}
				x := <-cx
				if x.ActiveIndex == sourceIndex {
					x.ActiveIndex = -1
					if x.SourceCompleted {
						sink(t)
					}
				}
				cx <- x
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
func Switch() rx.Operator {
	return SwitchMap(projectToObservable)
}

// SwitchMap creates an Observable that projects each source value to an
// Observable which is merged in the output Observable, emitting values only
// from the most recently projected Observable.
//
// SwitchMap maps each value to an Observable, then flattens all of these inner
// Observables using Switch.
func SwitchMap(project func(interface{}, int) (rx.Observable, error)) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := switchMapObservable{source, project}
		return rx.Create(obs.Subscribe)
	}
}

// SwitchMapTo creates an Observable that projects each source value to the
// same Observable which is flattened multiple times with Switch in the output
// Observable.
//
// It's like SwitchMap, but maps each value always to the same inner Observable.
func SwitchMapTo(inner rx.Observable) rx.Operator {
	return SwitchMap(func(interface{}, int) (rx.Observable, error) { return inner, nil })
}
