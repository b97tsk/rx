package rx

import (
	"context"
)

type congestingConcatObservable struct {
	Source  Observable
	Project func(interface{}, int) Observable
}

func (obs congestingConcatObservable) Subscribe(ctx context.Context, sink Observer) {
	var (
		sourceIndex = -1
		observer    Observer
	)

	observer = func(t Notification) {
		switch {
		case t.HasValue:
			sourceIndex++
			sourceIndex := sourceIndex
			sourceValue := t.Value

			obs := obs.Project(sourceValue, sourceIndex)

			childCtx, _ := obs.Subscribe(ctx, func(t Notification) {
				switch {
				case t.HasValue:
					sink(t)
				case t.HasError:
					observer = NopObserver
					sink(t)
				default:
					// do nothing
				}
			})

			<-childCtx.Done()

		default:
			sink(t)
		}
	}

	obs.Source.Subscribe(ctx, observer.Notify)
}

// CongestingConcat creates an output Observable which concurrently emits all
// values from every given input Observable.
//
// CongestingConcat flattens multiple Observables together by blending their
// values into one Observable.
//
// It's like Concat, but it congests the source.
func CongestingConcat(observables ...Observable) Observable {
	return FromObservables(observables...).Pipe(operators.CongestingConcatAll())
}

// CongestingConcatAll converts a higher-order Observable into a first-order
// Observable which concurrently delivers all values that are emitted on the
// inner Observables.
//
// It's like ConcatAll, but it congests the source.
func (Operators) CongestingConcatAll() Operator {
	return operators.CongestingConcatMap(ProjectToObservable)
}

// CongestingConcatMap creates an Observable that projects each source value to
// an Observable which is merged in the output Observable.
//
// CongestingConcatMap maps each value to an Observable, then flattens all of
// these inner Observables using CongestingConcatAll.
//
// It's like ConcatMap, but it congests the source.
func (Operators) CongestingConcatMap(project func(interface{}, int) Observable) Operator {
	return func(source Observable) Observable {
		obs := congestingConcatObservable{source, project}
		return Create(obs.Subscribe)
	}
}

// CongestingConcatMapTo creates an Observable that projects each source value
// to the same Observable which is merged multiple times in the output
// Observable.
//
// It's like CongestingConcatMap, but maps each value always to the same inner
// Observable.
//
// It's like ConcatMapTo, but it congests the source.
func (Operators) CongestingConcatMapTo(inner Observable) Operator {
	return operators.CongestingConcatMap(func(interface{}, int) Observable { return inner })
}
