package operators

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/queue"
)

// MergeScan applies an accumulator function over the source Observable where
// the accumulator function itself returns an Observable, then each
// intermediate Observable returned is merged into the output Observable.
//
// It's like Scan, but the Observables returned by the accumulator are merged
// into the outer Observable.
//
// For unlimited concurrency, passes -1.
func MergeScan(
	accumulator func(interface{}, interface{}, int) rx.Observable,
	seed interface{},
	concurrency int,
) rx.Operator {
	if accumulator == nil {
		panic("MergeScan: accumulator is nil")
	}

	if concurrency == 0 {
		concurrency = -1
	}

	return func(source rx.Observable) rx.Observable {
		return mergeScanObservable{
			Source:      source,
			Accumulator: accumulator,
			Seed:        seed,
			Concurrency: concurrency,
		}.Subscribe
	}
}

type mergeScanObservable struct {
	Source      rx.Observable
	Accumulator func(interface{}, interface{}, int) rx.Observable
	Seed        interface{}
	Concurrency int
}

type mergeScanState struct {
	Value interface{}
}

func (obs mergeScanObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel).Mutex()

	x := struct {
		sync.Mutex
		Queue     queue.Queue
		Index     int
		Workers   int
		State     atomic.Value
		Completed bool
	}{Index: -1}

	x.State.Store(mergeScanState{obs.Seed})

	var subscribeLocked func()

	subscribeLocked = func() {
		x.Index++

		sourceIndex := x.Index
		sourceValue := x.Queue.Pop()

		state := x.State.Load().(mergeScanState).Value

		obs1 := obs.Accumulator(state, sourceValue, sourceIndex)

		go obs1.Subscribe(ctx, func(t rx.Notification) {
			if t.HasValue {
				x.State.Store(mergeScanState{t.Value})
			}

			if t.HasValue || t.HasError {
				sink(t)

				return
			}

			x.Lock()
			defer x.Unlock()

			if x.Queue.Len() > 0 {
				subscribeLocked()
			} else {
				x.Workers--

				if x.Completed && x.Workers == 0 {
					sink(t)
				}
			}
		})
	}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		switch {
		case t.HasValue:
			x.Lock()
			defer x.Unlock()

			x.Queue.Push(t.Value)

			if x.Workers != obs.Concurrency {
				x.Workers++
				subscribeLocked()
			}

		case t.HasError:
			sink(t)

		default:
			x.Lock()
			defer x.Unlock()

			x.Completed = true

			if x.Workers == 0 {
				sink(t)
			}
		}
	})
}
