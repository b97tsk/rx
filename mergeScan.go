package rx

import (
	"context"

	"github.com/b97tsk/rx/x/queue"
)

// A MergeScanConfigure is a configure for MergeScan.
type MergeScanConfigure struct {
	Accumulator func(interface{}, interface{}) Observable
	Seed        interface{}
	Concurrent  int
}

// MakeFunc creates an OperatorFunc from this type.
func (conf MergeScanConfigure) MakeFunc() OperatorFunc {
	return MakeFunc(mergeScanOperator(conf).Call)
}

type mergeScanOperator MergeScanConfigure

func (op mergeScanOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Mutex(Finally(sink, cancel))

	type X struct {
		ActiveCount     int
		SourceCompleted bool
		Buffer          queue.Queue
		Seed            interface{}
		HasValue        bool
	}
	cx := make(chan *X, 1)
	cx <- &X{Seed: op.Seed}

	var doNextLocked func(*X)

	doNextLocked = func(x *X) {
		outerValue := x.Buffer.PopFront()

		// calls op.Accumulator synchronously
		obs := op.Accumulator(x.Seed, outerValue)

		go obs.Subscribe(ctx, func(t Notification) {
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
					x.ActiveCount--
					if x.ActiveCount == 0 && x.SourceCompleted {
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

	source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			x := <-cx
			x.Buffer.PushBack(t.Value)
			if x.ActiveCount != op.Concurrent {
				x.ActiveCount++
				doNextLocked(x)
			}
			cx <- x

		case t.HasError:
			sink(t)

		default:
			x := <-cx
			x.SourceCompleted = true
			if x.ActiveCount == 0 {
				if !x.HasValue {
					sink.Next(x.Seed)
				}
				sink(t)
			}
			cx <- x
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
		op := mergeScanOperator{accumulator, seed, -1}
		return source.Lift(op.Call)
	}
}
