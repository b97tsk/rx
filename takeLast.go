package rx

import (
	"container/ring"
	"context"
)

type takeLastOperator struct {
	Count int
}

func (op takeLastOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	var (
		head   *ring.Ring
		length int
	)
	return source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			if length < op.Count {
				tail := &ring.Ring{Value: t.Value}
				if head == nil {
					head = tail
				} else {
					tail.Link(head)
				}
				length++
			} else {
				head.Value = t.Value
				head = head.Next()
			}
		case t.HasError:
			sink(t)
		default:
			if length > 0 {
				done := ctx.Done()
				for i := 0; i < length; i++ {
					select {
					case <-done:
						return
					default:
					}
					sink.Next(head.Value)
					head = head.Next()
				}
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
func (Operators) TakeLast(count int) OperatorFunc {
	return func(source Observable) Observable {
		if count <= 0 {
			return Empty()
		}
		op := takeLastOperator{count}
		return source.Lift(op.Call)
	}
}
