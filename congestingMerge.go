package rx

import (
	"context"
	"sync"
)

type congestingMergeOperator struct {
	source     Operator
	project    func(interface{}, int) Observable
	concurrent int
}

func (op congestingMergeOperator) ApplyOptions(options []Option) Operator {
	for _, opt := range options {
		switch t := opt.(type) {
		case concurrentOption:
			op.concurrent = t.Value
		default:
			panic(ErrUnsupportedOption)
		}
	}
	return op
}

func (op congestingMergeOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()
	mu := sync.Mutex{}
	outerIndex := -1
	activeCount := 0
	completeSignal := make(chan struct{}, 1)

	concurrent := op.concurrent
	if concurrent == 0 {
		concurrent = -1
	}

	// Go statement makes this operator non-blocking.
	go op.source.Call(ctx, ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			mu.Lock()
			for activeCount == concurrent {
				mu.Unlock()
				select {
				case <-done:
					return
				case <-completeSignal:
				}
				mu.Lock()
			}
			defer mu.Unlock()

			activeCount++

			outerValue := t.Value
			outerIndex++
			outerIndex := outerIndex

			// calls project synchronously
			obsv := op.project(outerValue, outerIndex)

			go obsv.Subscribe(ctx, ObserverFunc(func(t Notification) {
				switch {
				case t.HasValue:
					ob.Next(t.Value)

				case t.HasError:
					ob.Error(t.Value.(error))
					cancel()

				default:
					mu.Lock()
					activeCount--
					mu.Unlock()

					select {
					case completeSignal <- struct{}{}:
					default:
					}
				}
			}))

		case t.HasError:
			ob.Error(t.Value.(error))
			cancel()

		default:
			mu.Lock()
			for activeCount > 0 {
				mu.Unlock()
				select {
				case <-done:
					return
				case <-completeSignal:
				}
				mu.Lock()
			}
			mu.Unlock()
			ob.Complete()
			cancel()
		}
	}))

	return ctx, cancel
}

// CongestingMerge creates an output Observable which concurrently emits all
// values from every given input Observable.
//
// CongestingMerge flattens multiple Observables together by blending their
// values into one Observable.
//
// It's like Merge, but it may congest the source due to concurrent limit.
func CongestingMerge(observables ...interface{}) Observable {
	return FromSlice(observables).CongestingMergeAll()
}

// CongestingMergeAll converts a higher-order Observable into a first-order
// Observable which concurrently delivers all values that are emitted on the
// inner Observables.
//
// It's like MergeAll, but it may congest the source due to concurrent limit.
func (o Observable) CongestingMergeAll() Observable {
	op := congestingMergeOperator{
		source:     o.Op,
		project:    projectToObservable,
		concurrent: -1,
	}
	return Observable{op}.Mutex()
}

// CongestingMergeMap creates an Observable that projects each source value to
// an Observable which is merged in the output Observable.
//
// CongestingMergeMap maps each value to an Observable, then flattens all of
// these inner Observables using CongestingMergeAll.
//
// It's like MergeMap, but it may congest the source due to concurrent limit.
func (o Observable) CongestingMergeMap(project func(interface{}, int) Observable) Observable {
	op := congestingMergeOperator{
		source:     o.Op,
		project:    project,
		concurrent: -1,
	}
	return Observable{op}.Mutex()
}

// CongestingMergeMapTo creates an Observable that projects each source value
// to the same Observable which is merged multiple times in the output
// Observable.
//
// It's like CongestingMergeMap, but maps each value always to the same inner
// Observable.
//
// It's like MergeMapTo, but it may congest the source due to concurrent limit.
func (o Observable) CongestingMergeMapTo(inner Observable) Observable {
	op := congestingMergeOperator{
		source:     o.Op,
		project:    func(interface{}, int) Observable { return inner },
		concurrent: -1,
	}
	return Observable{op}.Mutex()
}
