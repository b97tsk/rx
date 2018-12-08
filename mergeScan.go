package rx

import (
	"container/list"
	"context"
	"sync"
)

// MergeScanOperator is an operator type.
type MergeScanOperator struct {
	Accumulator func(interface{}, interface{}) Observable
	Seed        interface{}
	Concurrent  int
}

// MakeFunc creates an OperatorFunc from this operator.
func (op MergeScanOperator) MakeFunc() OperatorFunc {
	return MakeFunc(op.Call)
}

// Call invokes an execution of this operator.
func (op MergeScanOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Mutex(Finally(sink, cancel))

	var (
		mutex           sync.Mutex
		activeCount     int
		sourceCompleted bool
		seed            = op.Seed
		hasValue        bool
		buffer          list.List
		doNextLocked    func()
	)

	concurrent := op.Concurrent
	if concurrent == 0 {
		concurrent = -1
	}

	doNextLocked = func() {
		outerValue := buffer.Remove(buffer.Front())

		// calls op.Accumulator synchronously
		obs := op.Accumulator(seed, outerValue)

		go obs.Subscribe(ctx, func(t Notification) {
			switch {
			case t.HasValue:
				mutex.Lock()
				seed = t.Value
				hasValue = true
				mutex.Unlock()

				sink(t)

			case t.HasError:
				sink(t)

			default:
				mutex.Lock()
				if buffer.Len() > 0 {
					doNextLocked()
				} else {
					activeCount--
					if activeCount == 0 && sourceCompleted {
						if !hasValue {
							sink.Next(seed)
						}
						sink(t)
					}
				}
				mutex.Unlock()
			}
		})
	}

	source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			mutex.Lock()
			buffer.PushBack(t.Value)
			if activeCount != concurrent {
				activeCount++
				doNextLocked()
			}
			mutex.Unlock()

		case t.HasError:
			sink(t)

		default:
			mutex.Lock()
			sourceCompleted = true
			if activeCount == 0 {
				if !hasValue {
					sink.Next(seed)
				}
				sink(t)
			}
			mutex.Unlock()
		}
	})

	return ctx, cancel
}

// MergeScan applies an accumulator function over the source Observable where
// the accumulator function itself returns an Observable, then each
// intermediate Observable returned is merged into the output Observable.
//
// It's like Scan, but the Observables returned by the accumulator are merged
// into the outer Observable.
func (Operators) MergeScan(accumulator func(interface{}, interface{}) Observable, seed interface{}) OperatorFunc {
	return func(source Observable) Observable {
		op := MergeScanOperator{accumulator, seed, -1}
		return source.Lift(op.Call)
	}
}
