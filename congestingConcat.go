package rx

import (
	"context"
)

type congestingConcatOperator struct {
	project func(interface{}, int) Observable
}

func (op congestingConcatOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	var (
		outerIndex      = -1
		mutableObserver Observer
	)

	mutableObserver = func(t Notification) {
		switch {
		case t.HasValue:
			outerValue := t.Value
			outerIndex++
			outerIndex := outerIndex

			obsv := op.project(outerValue, outerIndex)

			childCtx, _ := obsv.Subscribe(ctx, func(t Notification) {
				switch {
				case t.HasValue:
					ob.Next(t.Value)
				case t.HasError:
					mutableObserver = NopObserver
					ob.Error(t.Value.(error))
					cancel()
				default:
				}
			})

			<-childCtx.Done()

		case t.HasError:
			ob.Error(t.Value.(error))
			cancel()

		default:
			ob.Complete()
			cancel()
		}
	}

	source.Subscribe(ctx, func(t Notification) { t.Observe(mutableObserver) })

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
	op := congestingConcatOperator{projectToObservable}
	return o.Lift(op.Call)
}

// CongestingConcatMap creates an Observable that projects each source value to
// an Observable which is merged in the output Observable.
//
// CongestingConcatMap maps each value to an Observable, then flattens all of
// these inner Observables using CongestingConcatAll.
//
// It's like ConcatMap, but it congests the source.
func (o Observable) CongestingConcatMap(project func(interface{}, int) Observable) Observable {
	op := congestingConcatOperator{project}
	return o.Lift(op.Call)
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
	return o.CongestingConcatMap(func(interface{}, int) Observable { return inner })
}
