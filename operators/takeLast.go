package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/queue"
)

type takeLastObservable struct {
	Source rx.Observable
	Count  int
}

func (obs takeLastObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	var queue queue.Queue
	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		switch {
		case t.HasValue:
			if queue.Len() == obs.Count {
				queue.Pop()
			}
			queue.Push(t.Value)
		case t.HasError:
			sink(t)
		default:
			for i, j := 0, queue.Len(); i < j; i++ {
				if ctx.Err() != nil {
					return
				}
				sink.Next(queue.At(i))
			}
			sink(t)
		}
	})
}

// TakeLast creates an Observable that emits only the last count values emitted
// by the source Observable.
//
// TakeLast remembers the latest count values, then emits those only when the
// source completes.
func TakeLast(count int) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		if count <= 0 {
			return rx.Empty()
		}
		return takeLastObservable{source, count}.Subscribe
	}
}
