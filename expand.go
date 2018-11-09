package rx

import (
	"container/list"
	"context"
	"sync"
)

// ExpandOperator is an operator type.
type ExpandOperator struct {
	Project    func(interface{}, int) Observable
	Concurrent int
}

// MakeFunc creates an OperatorFunc from this operator.
func (op ExpandOperator) MakeFunc() OperatorFunc {
	return MakeFunc(op.Call)
}

// Call invokes an execution of this operator.
func (op ExpandOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Mutex(Finally(sink, cancel))

	var (
		mutex           sync.Mutex
		outerIndex      = -1
		activeCount     int
		sourceCompleted bool
		buffer          list.List
		doNextLocked    func()
	)

	concurrent := op.Concurrent
	if concurrent == 0 {
		concurrent = -1
	}

	doNextLocked = func() {
		outerIndex++
		outerIndex := outerIndex
		outerValue := buffer.Remove(buffer.Front())

		sink.Next(outerValue)

		// calls op.Project synchronously
		obsv := op.Project(outerValue, outerIndex)

		go obsv.Subscribe(ctx, func(t Notification) {
			switch {
			case t.HasValue:
				mutex.Lock()
				defer mutex.Unlock()

				buffer.PushBack(t.Value)

				if activeCount != concurrent {
					activeCount++
					doNextLocked()
				}

			case t.HasError:
				sink(t)

			default:
				mutex.Lock()
				defer mutex.Unlock()

				if buffer.Len() > 0 {
					doNextLocked()
					break
				}

				activeCount--

				if activeCount == 0 && sourceCompleted {
					sink(t)
				}
			}
		})
	}

	source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			mutex.Lock()
			defer mutex.Unlock()

			buffer.PushBack(t.Value)

			if activeCount != concurrent {
				activeCount++
				doNextLocked()
			}

		case t.HasError:
			sink(t)

		default:
			mutex.Lock()
			defer mutex.Unlock()
			sourceCompleted = true
			if activeCount == 0 {
				sink(t)
			}
		}
	})

	return ctx, cancel
}

// Expand recursively projects each source value to an Observable which is
// merged in the output Observable.
//
// It's similar to MergeMap, but applies the projection function to every
// source value as well as every output value. It's recursive.
func (Operators) Expand(project func(interface{}, int) Observable) OperatorFunc {
	return func(source Observable) Observable {
		op := ExpandOperator{project, -1}
		return source.Lift(op.Call)
	}
}
