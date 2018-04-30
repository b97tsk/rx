package rx

import (
	"context"
)

type concatMapOperator mergeMapOperator

func (op concatMapOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	return mergeMapOperator(op).Call(ctx, ob)
}

// Concat creates an output Observable which sequentially emits all values
// from given Observable and then moves on to the next.
//
// Concat concatenates multiple Observables together by sequentially emitting
// their values, one Observable after the other.
func Concat(observables ...interface{}) Observable {
	return FromSlice(observables).ConcatAll()
}

// ConcatAll converts a higher-order Observable into a first-order Observable
// by concatenating the inner Observables in order.
//
// ConcatAll flattens an Observable-of-Observables by putting one inner
// Observable after the other.
func (o Observable) ConcatAll() Observable {
	op := concatMapOperator{
		source:     o.Op,
		project:    projectToObservable,
		concurrent: 1,
	}
	return Observable{op}
}

// ConcatMap projects each source value to an Observable which is merged in
// the output Observable, in a serialized fashion waiting for each one to
// complete before merging the next.
//
// ConcatMap maps each value to an Observable, then flattens all of these inner
// Observables using ConcatAll.
func (o Observable) ConcatMap(project func(interface{}, int) Observable) Observable {
	op := concatMapOperator{
		source:     o.Op,
		project:    project,
		concurrent: 1,
	}
	return Observable{op}
}

// ConcatMapTo projects each source value to the same Observable which is
// merged multiple times in a serialized fashion on the output Observable.
//
// It's like ConcatMap, but maps each value always to the same inner Observable.
func (o Observable) ConcatMapTo(inner Observable) Observable {
	op := concatMapOperator{
		source:     o.Op,
		project:    func(interface{}, int) Observable { return inner },
		concurrent: 1,
	}
	return Observable{op}
}
