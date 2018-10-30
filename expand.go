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
	done := ctx.Done()
	sink = Mutex(sink)

	var (
		mu             sync.Mutex
		outerIndex     = -1
		activeCount    = 0
		buffer         list.List
		completeSignal = make(chan struct{}, 1)
		doNextLocked   func()
	)

	concurrent := op.Concurrent
	if concurrent == 0 {
		concurrent = -1
	}

	doNextLocked = func() {
		outerValue := buffer.Remove(buffer.Front())
		outerIndex++
		outerIndex := outerIndex

		sink.Next(outerValue)

		// calls op.Project synchronously
		obsv := op.Project(outerValue, outerIndex)

		go obsv.Subscribe(ctx, func(t Notification) {
			switch {
			case t.HasValue:
				mu.Lock()
				defer mu.Unlock()

				buffer.PushBack(t.Value)

				if activeCount != concurrent {
					activeCount++
					doNextLocked()
				}

			case t.HasError:
				sink(t)
				cancel()

			default:
				mu.Lock()

				if buffer.Len() > 0 {
					defer mu.Unlock()
					doNextLocked()
					break
				}

				activeCount--
				mu.Unlock()

				select {
				case completeSignal <- struct{}{}:
				default:
				}
			}
		})
	}

	source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			mu.Lock()
			defer mu.Unlock()

			buffer.PushBack(t.Value)

			if activeCount != concurrent {
				activeCount++
				doNextLocked()
			}

		case t.HasError:
			sink(t)
			cancel()

		default:
			mu.Lock()
			if activeCount > 0 {
				go func() {
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
					sink.Complete()
					cancel()
				}()
				return
			}
			mu.Unlock()
			sink(t)
			cancel()
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
