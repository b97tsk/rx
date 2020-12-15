package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// ConcatSyncAll creates an Observable that flattens a higher-order
// Observable into a first-order Observable by concatenating the inner
// Observables in order.
//
// It's like ConcatAll, but it does not buffer the source.
func ConcatSyncAll() rx.Operator {
	return ConcatSyncMap(projectToObservable)
}

// ConcatSyncMap creates an Observable that converts the source Observable
// into a higher-order Observable, by projecting each source value to an
// Observable, and flattens it into a first-order Observable by concatenating
// the inner Observables in order.
//
// It's like ConcatMap, but it does not buffer the source.
func ConcatSyncMap(project func(interface{}, int) rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return concatSyncObservable{source, project}.Subscribe
	}
}

// ConcatSyncMapTo creates an Observable that converts the source Observable
// into a higher-order Observable, by projecting each source value to the same
// Observable, and flattens it into a first-order Observable by concatenating
// the inner Observables in order.
//
// It's like ConcatSyncMap, but maps each value always to the same inner
// Observable.
//
// It's like ConcatMapTo, but it does not buffer the source.
func ConcatSyncMapTo(inner rx.Observable) rx.Operator {
	return ConcatSyncMap(func(interface{}, int) rx.Observable { return inner })
}

type concatSyncObservable struct {
	Source  rx.Observable
	Project func(interface{}, int) rx.Observable
}

func (obs concatSyncObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	var observer rx.Observer

	sourceIndex := -1

	observer = func(t rx.Notification) {
		switch {
		case t.HasValue:
			sourceIndex++

			obs1 := obs.Project(t.Value, sourceIndex)

			err := obs1.BlockingSubscribe(ctx, func(t rx.Notification) {
				if t.HasValue || t.HasError {
					sink(t)
				}
			})
			if err != nil {
				observer = rx.Noop
			}

		default:
			sink(t)
		}
	}

	obs.Source.Subscribe(ctx, observer.Sink)
}
