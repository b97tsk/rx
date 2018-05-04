package rx

import (
	"context"
)

type congestingConcatOperator struct {
	source  Operator
	project func(interface{}, int) Observable
}

func (op congestingConcatOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	outerIndex := -1

	mutable := MutableObserver{}

	mutable.Observer = ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			outerValue := t.Value
			outerIndex++
			outerIndex := outerIndex

			obsv := op.project(outerValue, outerIndex)

			childCtx, _ := obsv.Subscribe(ctx, ObserverFunc(func(t Notification) {
				switch {
				case t.HasValue:
					ob.Next(t.Value)
				case t.HasError:
					mutable.Observer = NopObserver
					ob.Error(t.Value.(error))
					cancel()
				default:
				}
			}))

			<-childCtx.Done()

		case t.HasError:
			ob.Error(t.Value.(error))
			cancel()

		default:
			ob.Complete()
			cancel()
		}
	})

	// Go statement makes this operator non-blocking.
	go op.source.Call(ctx, &mutable)

	return ctx, cancel
}

// CongestingConcat creates an output Observable which concurrently emits all
// values from every given input Observable.
//
// CongestingConcat flattens multiple Observables together by blending their
// values into one Observable.
//
// It's like Concat, but it congests the source.
func CongestingConcat(observables ...interface{}) Observable {
	return FromSlice(observables).CongestingConcatAll()
}

// CongestingConcatAll converts a higher-order Observable into a first-order
// Observable which concurrently delivers all values that are emitted on the
// inner Observables.
//
// It's like ConcatAll, but it congests the source.
func (o Observable) CongestingConcatAll() Observable {
	op := congestingConcatOperator{
		source:  o.Op,
		project: projectToObservable,
	}
	return Observable{op}
}

// CongestingConcatMap creates an Observable that projects each source value to
// an Observable which is merged in the output Observable.
//
// CongestingConcatMap maps each value to an Observable, then flattens all of
// these inner Observables using CongestingConcatAll.
//
// It's like ConcatMap, but it congests the source.
func (o Observable) CongestingConcatMap(project func(interface{}, int) Observable) Observable {
	op := congestingConcatOperator{
		source:  o.Op,
		project: project,
	}
	return Observable{op}
}

// CongestingConcatMapTo creates an Observable that projects each source value
// to the same Observable which is merged multiple times in the output
// Observable.
//
// It's like CongestingConcatMap, but maps each value always to the same inner
// Observable.
//
// It's like ConcatMapTo, but it congests the source.
func (o Observable) CongestingConcatMapTo(inner Observable) Observable {
	op := congestingConcatOperator{
		source:  o.Op,
		project: func(interface{}, int) Observable { return inner },
	}
	return Observable{op}
}
