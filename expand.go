package rx

import (
	"context"
	"sync"

	"github.com/b97tsk/rx/x/queue"
)

// An ExpandConfigure is a configure for Expand.
type ExpandConfigure struct {
	Project    func(interface{}, int) Observable
	Concurrent int
}

// MakeFunc creates an OperatorFunc from this type.
func (conf ExpandConfigure) MakeFunc() OperatorFunc {
	return MakeFunc(expandOperator(conf).Call)
}

type expandOperator ExpandConfigure

func (op expandOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Mutex(Finally(sink, cancel))

	var (
		mutex           sync.Mutex
		outerIndex      = -1
		activeCount     int
		sourceCompleted bool
		buffer          queue.Queue
		doNextLocked    func()
	)

	doNextLocked = func() {
		outerIndex++
		outerIndex := outerIndex
		outerValue := buffer.PopFront()

		sink.Next(outerValue)

		// calls op.Project synchronously
		obs := op.Project(outerValue, outerIndex)

		go obs.Subscribe(ctx, func(t Notification) {
			switch {
			case t.HasValue:
				mutex.Lock()
				buffer.PushBack(t.Value)
				if activeCount != op.Concurrent {
					activeCount++
					doNextLocked()
				}
				mutex.Unlock()

			case t.HasError:
				sink(t)

			default:
				mutex.Lock()
				if buffer.Len() > 0 {
					doNextLocked()
				} else {
					activeCount--
					if activeCount == 0 && sourceCompleted {
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
			if activeCount != op.Concurrent {
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
				sink(t)
			}
			mutex.Unlock()
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
		op := expandOperator{project, -1}
		return source.Lift(op.Call)
	}
}
