package rx

import (
	"context"
)

type congestOperator struct {
	source   Operator
	capacity int
}

func (op congestOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()
	c := make(chan Notification, op.capacity)

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-c:
				switch {
				case t.HasValue:
					ob.Next(t.Value)
				case t.HasError:
					ob.Error(t.Value.(error))
					cancel()
					return
				default:
					ob.Complete()
					cancel()
					return
				}
			}
		}
	}()

	op.source.Call(ctx, func(t Notification) {
		select {
		case <-done:
		case c <- t:
		}
	})

	return ctx, cancel
}

// Congest creates an Observable that mirrors the source Observable, caches
// emissions if the source emits too fast, and congests the source if the cache
// is full.
func (o Observable) Congest(capacity int) Observable {
	if capacity < 1 {
		return o
	}
	op := congestOperator{
		source:   o.Op,
		capacity: capacity,
	}
	return Observable{op}
}
