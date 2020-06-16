package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

type flatObservable struct {
	Source  rx.Observable
	Flat    func(observables ...rx.Observable) rx.Observable
	Project func(interface{}, int) (rx.Observable, error)
}

func (obs flatObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	var (
		observables []rx.Observable
		observer    rx.Observer
	)

	sourceIndex := -1

	observer = func(t rx.Notification) {
		switch {
		case t.HasValue:
			sourceIndex++

			obs, err := obs.Project(t.Value, sourceIndex)
			if err != nil {
				observer = rx.Noop
				sink.Error(err)
				return
			}

			observables = append(observables, obs)

		case t.HasError:
			sink(t)

		default:
			obs := obs.Flat(observables...)
			obs.Subscribe(ctx, sink)
		}
	}

	obs.Source.Subscribe(ctx, observer.Sink)
}

// Flat creates an Observable that flattens a higher-order Observable into
// a first-order Observable, by applying a flat function to those Observables
// emitted by the higher-order Observable, and starts subscribing to it.
func Flat(flat func(observables ...rx.Observable) rx.Observable) rx.Operator {
	return FlatMap(flat, func(val interface{}, idx int) (rx.Observable, error) {
		if obs, ok := val.(rx.Observable); ok {
			return obs, nil
		}
		return nil, rx.ErrNotObservable
	})
}

// FlatMap creates an Observable that converts the source Observable to a
// higher-order Observable, by projecting each source value to an Observable,
// and flattens it into a first-order Observable, by applying a flat function
// to those Observables emitted by the higher-order Observable, and starts
// subscribing to it.
func FlatMap(flat func(observables ...rx.Observable) rx.Observable, project func(interface{}, int) (rx.Observable, error)) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := flatObservable{source, flat, project}
		return rx.Create(obs.Subscribe)
	}
}
