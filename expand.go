package rx

import (
	"context"

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

	type X struct {
		Index           int
		ActiveCount     int
		SourceCompleted bool
		Buffer          queue.Queue
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	var doNextLocked func(*X)

	doNextLocked = func(x *X) {
		outerIndex := x.Index
		outerValue := x.Buffer.PopFront()
		x.Index++

		sink.Next(outerValue)

		// calls op.Project synchronously
		obs := op.Project(outerValue, outerIndex)

		go obs.Subscribe(ctx, func(t Notification) {
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
				if x.Buffer.Len() > 0 {
					doNextLocked(x)
				} else {
					x.ActiveCount--
					if x.ActiveCount == 0 && x.SourceCompleted {
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
				sink(t)
			}
			cx <- x
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
