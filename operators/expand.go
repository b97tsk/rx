package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/queue"
)

// An ExpandConfigure is a configure for Expand.
type ExpandConfigure struct {
	Project    func(interface{}) rx.Observable
	Concurrent int
}

// Make creates an Operator from this configure.
func (configure ExpandConfigure) Make() rx.Operator {
	if configure.Project == nil {
		panic("Expand: nil Project")
	}
	if configure.Concurrent == 0 {
		configure.Concurrent = -1
	}
	return func(source rx.Observable) rx.Observable {
		return expandObservable{source, configure}.Subscribe
	}
}

type expandObservable struct {
	Source rx.Observable
	ExpandConfigure
}

func (obs expandObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)
	sink = sink.WithCancel(cancel).Mutex()

	type X struct {
		Active          int
		Buffer          queue.Queue
		SourceCompleted bool
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	var doNextLocked func(*X)

	doNextLocked = func(x *X) {
		val := x.Buffer.Pop()

		sink.Next(val)

		obs1 := obs.Project(val)

		go obs1.Subscribe(ctx, func(t rx.Notification) {
			switch {
			case t.HasValue:
				x := <-cx
				x.Buffer.Push(t.Value)
				if x.Active != obs.Concurrent {
					x.Active++
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
					x.Active--
					if x.Active == 0 && x.SourceCompleted {
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
			x.Buffer.Push(t.Value)
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
				sink(t)
			}
			cx <- x
		}
	})
}

// Expand recursively projects each source value to an Observable which is
// merged in the output Observable.
//
// It's similar to MergeMap, but applies the projection function to every
// source value as well as every output value. It's recursive.
func Expand(project func(interface{}) rx.Observable) rx.Operator {
	return ExpandConfigure{project, -1}.Make()
}
