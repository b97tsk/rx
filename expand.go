package rx

import (
	"container/list"
	"context"
	"sync"
)

type expandOperator struct {
	source     Operator
	project    func(interface{}, int) Observable
	concurrent int
}

func (op expandOperator) ApplyOptions(options []Option) Operator {
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

func (op expandOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()
	mu := sync.Mutex{}
	outerIndex := -1
	activeCount := 0
	buffer := list.List{}
	completeSignal := make(chan struct{}, 1)
	ob = Normalize(ob)

	concurrent := op.concurrent
	if concurrent == 0 {
		concurrent = -1
	}

	var doNextLocked func()

	doNextLocked = func() {
		outerValue := buffer.Remove(buffer.Front())
		outerIndex++
		outerIndex := outerIndex

		ob.Next(outerValue)

		// calls project synchronously
		obsv := op.project(outerValue, outerIndex)

		go obsv.Subscribe(ctx, ObserverFunc(func(t Notification) {
			switch {
			case t.HasValue:
				mu.Lock()

				buffer.PushBack(t.Value)

				if activeCount != concurrent {
					activeCount++
					doNextLocked()
				}

				mu.Unlock()

			case t.HasError:
				ob.Error(t.Value.(error))
				cancel()

			default:
				mu.Lock()

				if buffer.Len() > 0 {
					doNextLocked()
					mu.Unlock()
					break
				}

				activeCount--
				mu.Unlock()

				select {
				case completeSignal <- struct{}{}:
				default:
				}
			}
		}))
	}

	// Go statement makes this operator non-blocking.
	go op.source.Call(ctx, ObserverFunc(func(t Notification) {
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

// Expand recursively projects each source value to an Observable which is
// merged in the output Observable.
//
// It's similar to MergeMap, but applies the projection function to every
// source value as well as every output value. It's recursive.
func (o Observable) Expand(project func(interface{}, int) Observable) Observable {
	op := expandOperator{
		source:     o.Op,
		project:    project,
		concurrent: -1,
	}
	return Observable{op}
}
