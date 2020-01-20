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

// Use creates an Operator from this configure.
func (configure ExpandConfigure) Use() Operator {
	return func(source Observable) Observable {
		return expandObservable{source, configure}.Subscribe
	}
}

type expandObservable struct {
	Source Observable
	ExpandConfigure
}

func (obs expandObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
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
		sourceIndex := x.Index
		sourceValue := x.Buffer.PopFront()
		x.Index++

		sink.Next(sourceValue)

		// calls obs.Project synchronously
		obs1 := obs.Project(sourceValue, sourceIndex)

		go obs1.Subscribe(ctx, func(t Notification) {
			switch {
			case t.HasValue:
				x := <-cx
				x.Buffer.PushBack(t.Value)
				if x.ActiveCount != obs.Concurrent {
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

	obs.Source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			x := <-cx
			x.Buffer.PushBack(t.Value)
			if x.ActiveCount != obs.Concurrent {
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
func (Operators) Expand(project func(interface{}, int) Observable) Operator {
	return ExpandConfigure{project, -1}.Use()
}
