package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/queue"
)

// A MergeScanConfigure is a configure for MergeScan.
type MergeScanConfigure struct {
	Accumulator func(interface{}, interface{}) rx.Observable
	Seed        interface{}
	Concurrent  int
}

// Use creates an Operator from this configure.
func (configure MergeScanConfigure) Use() rx.Operator {
	if configure.Accumulator == nil {
		panic("MergeScan: nil Accumulator")
	}
	if configure.Concurrent == 0 {
		configure.Concurrent = -1
	}
	return func(source rx.Observable) rx.Observable {
		obs := mergeScanObservable{source, configure}
		return rx.Create(obs.Subscribe)
	}
}

type mergeScanObservable struct {
	Source rx.Observable
	MergeScanConfigure
}

func (obs mergeScanObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	sink = rx.Mutex(sink)

	type X struct {
		Active          int
		Buffer          queue.Queue
		Seed            interface{}
		HasValue        bool
		SourceCompleted bool
	}
	cx := make(chan *X, 1)
	cx <- &X{Seed: obs.Seed}

	var doNextLocked func(*X)

	doNextLocked = func(x *X) {
		sourceValue := x.Buffer.PopFront()

		// Call obs.Accumulator synchronously.
		obs := obs.Accumulator(x.Seed, sourceValue)

		go obs.Subscribe(ctx, func(t rx.Notification) {
			switch {
			case t.HasValue:
				x := <-cx
				x.Seed = t.Value
				x.HasValue = true
				cx <- x

				sink(t)

			case t.HasError:
				sink(t)

			default:
				x := <-cx
				if x.Buffer.Len() > 0 {
					doNextLocked(x)
				} else {
					x.Active--
					if x.Active == 0 && x.SourceCompleted {
						if !x.HasValue {
							sink.Next(x.Seed)
						}
						sink(t)
					}
				}
				cx <- x
			}
		})
	}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		switch {
		case t.HasValue:
			x := <-cx
			x.Buffer.PushBack(t.Value)
			if x.Active != obs.Concurrent {
				x.Active++
				doNextLocked(x)
			}
			cx <- x

		case t.HasError:
			sink(t)

		default:
			x := <-cx
			x.SourceCompleted = true
			if x.Active == 0 {
				if !x.HasValue {
					sink.Next(x.Seed)
				}
				sink(t)
			}
			cx <- x
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
	return MergeScanConfigure{accumulator, seed, -1}.Use()
}
