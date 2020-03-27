package rx

import (
	"context"
)

// A ToObservablesConfigure is a configure for ToObservables.
type ToObservablesConfigure struct {
	Flat func(observables ...Observable) Observable
}

// Use creates an Operator from this configure.
func (configure ToObservablesConfigure) Use() Operator {
	return func(source Observable) Observable {
		obs := toObservablesObservable{source, configure}
		return Create(obs.Subscribe)
	}
}

type toObservablesObservable struct {
	Source Observable
	ToObservablesConfigure
}

func (obs toObservablesObservable) Subscribe(ctx context.Context, sink Observer) {
	var (
		observables []Observable
		observer    Observer
	)

	observer = func(t Notification) {
		switch {
		case t.HasValue:
			if obs, ok := t.Value.(Observable); ok {
				observables = append(observables, obs)
			} else {
				observer = NopObserver
				sink.Error(ErrNotObservable)
			}
		case t.HasError:
			sink(t)
		default:
			if obs.Flat != nil {
				obs := obs.Flat(observables...)
				obs.Subscribe(ctx, sink)
			} else {
				sink.Next(observables)
				sink.Complete()
			}
		}
	}

	obs.Source.Subscribe(ctx, observer.Notify)
}

// ToObservables creates an Observable that collects all the Observables the
// source emits, then emits them as a slice of Observable when the source
// completes.
func (Operators) ToObservables() Operator {
	return ToObservablesConfigure{}.Use()
}
