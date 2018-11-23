package rx

import (
	"context"
)

// ToObservablesOperator is an operator type.
type ToObservablesOperator struct {
	Flat func(observables ...Observable) Observable
}

// MakeFunc creates an OperatorFunc from this operator.
func (op ToObservablesOperator) MakeFunc() OperatorFunc {
	return MakeFunc(op.Call)
}

// Call invokes an execution of this operator.
func (op ToObservablesOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

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
			if op.Flat != nil {
				obs := op.Flat(observables...)
				obs.Subscribe(ctx, sink)
			} else {
				sink.Next(observables)
				sink.Complete()
			}
		}
	}

	source.Subscribe(ctx, observer.Notify)

	return ctx, cancel
}

// ToObservables creates an Observable that collects all the Observables the
// source emits, then emits them as a slice of Observable when the source
// completes.
func (Operators) ToObservables() OperatorFunc {
	return func(source Observable) Observable {
		op := ToObservablesOperator{}
		return source.Lift(op.Call)
	}
}
