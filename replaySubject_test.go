package rx_test

import (
	"context"
	"testing"
	"time"

	. "github.com/b97tsk/rx"
)

func TestReplaySubject(t *testing.T) {
	t.Run("#1", func(t *testing.T) {
		subject := NewReplaySubject(3, 0)
		subscribeThenComplete := Observable(
			func(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(ctx)
				defer cancel()
				sink = Mutex(sink)
				subject.Subscribe(ctx, sink)
				sink.Complete()
				return ctx, cancel
			},
		)
		subject.Next("A")
		subscribe(t, []Observable{subscribeThenComplete}, "A", xComplete)
		subject.Next("B")
		subscribe(t, []Observable{subscribeThenComplete}, "A", "B", xComplete)
		subject.Next("C")
		subscribe(t, []Observable{subscribeThenComplete}, "A", "B", "C", xComplete)
		subject.Next("D")
		subscribe(t, []Observable{subscribeThenComplete}, "B", "C", "D", xComplete)
		subject.Error(xErrTest)
		subscribe(t, []Observable{subscribeThenComplete}, xErrTest)
	})
	t.Run("#2", func(t *testing.T) {
		subject := NewReplaySubject(0, step(5))
		subscribeThenComplete := Observable(
			func(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(ctx)
				defer cancel()
				sink = Mutex(sink)
				subject.Subscribe(ctx, sink)
				sink.Complete()
				return ctx, cancel
			},
		)
		subject.Next("A")
		subscribe(t, []Observable{subscribeThenComplete}, "A", xComplete)
		time.Sleep(step(2))
		subject.Next("B")
		subscribe(t, []Observable{subscribeThenComplete}, "A", "B", xComplete)
		time.Sleep(step(2))
		subject.Next("C")
		subscribe(t, []Observable{subscribeThenComplete}, "A", "B", "C", xComplete)
		time.Sleep(step(2))
		subject.Next("D")
		subject.Complete()
		subscribe(t, []Observable{subscribeThenComplete}, "B", "C", "D", xComplete)
		time.Sleep(step(2))
		subscribe(t, []Observable{subscribeThenComplete}, "C", "D", xComplete)
		time.Sleep(step(2))
		subscribe(t, []Observable{subscribeThenComplete}, "D", xComplete)
		time.Sleep(step(2))
		subscribe(t, []Observable{subscribeThenComplete}, xComplete)
	})
}
