package operators

import (
	"context"
	"sync"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/queue"
)

// A MergeScanConfigure is a configure for MergeScan.
type MergeScanConfigure struct {
	Accumulator func(interface{}, interface{}) rx.Observable
	Seed        interface{}
	Concurrency int
}

// Make creates an Operator from this configure.
func (configure MergeScanConfigure) Make() rx.Operator {
	if configure.Accumulator == nil {
		panic("MergeScan: Accumulator is nil")
	}
	if configure.Concurrency == 0 {
		configure.Concurrency = -1
	}
	return func(source rx.Observable) rx.Observable {
		return mergeScanObservable{source, configure}.Subscribe
	}
}

type mergeScanObservable struct {
	Source rx.Observable
	MergeScanConfigure
}

func (obs mergeScanObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)
	sink = sink.WithCancel(cancel).Mutex()

	x := struct {
		sync.Mutex
		Queue     queue.Queue
		Workers   int
		Seed      interface{}
		HasValue  bool
		Completed bool
	}{Seed: obs.Seed}

	var subscribeLocked func()
	subscribeLocked = func() {
		sourceValue := x.Queue.Pop()
		obs1 := obs.Accumulator(x.Seed, sourceValue)
		go obs1.Subscribe(ctx, func(t rx.Notification) {
			switch {
			case t.HasValue:
				x.Lock()
				x.Seed = t.Value
				x.HasValue = true
				x.Unlock()

				sink(t)

			case t.HasError:
				sink(t)

			default:
				x.Lock()
				defer x.Unlock()
				if x.Queue.Len() > 0 {
					subscribeLocked()
				} else {
					x.Workers--
					if x.Completed && x.Workers == 0 {
						if !x.HasValue {
							sink.Next(x.Seed)
						}
						sink(t)
					}
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
				if !x.HasValue {
					sink.Next(x.Seed)
				}
				sink(t)
			}
		}
	})
}

// MergeScan applies an accumulator function over the source Observable where
// the accumulator function itself returns an Observable, then each
// intermediate Observable returned is merged into the output Observable.
//
// It's like Scan, but the Observables returned by the accumulator are merged
// into the outer Observable.
func MergeScan(accumulator func(interface{}, interface{}) rx.Observable, seed interface{}) rx.Operator {
	return MergeScanConfigure{accumulator, seed, -1}.Make()
}
