package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// A ToObservablesConfigure is a configure for ToObservables.
type ToObservablesConfigure struct {
	Flat func(observables ...rx.Observable) rx.Observable
}

// Use creates an Operator from this configure.
func (configure ToObservablesConfigure) Use() rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := toObservablesObservable{source, configure}
		return rx.Create(obs.Subscribe)
	}
}

type toObservablesObservable struct {
	Source rx.Observable
	ToObservablesConfigure
}

func (obs toObservablesObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	var (
		observables []rx.Observable
		observer    rx.Observer
	)

	observer = func(t rx.Notification) {
		switch {
		case t.HasValue:
			if obs, ok := t.Value.(rx.Observable); ok {
				observables = append(observables, obs)
			} else {
				observer = rx.Noop
				sink.Error(rx.ErrNotObservable)
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
func ToObservables() rx.Operator {
	return ToObservablesConfigure{}.Use()
}
