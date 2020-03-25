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
		subscribe(t, []Observable{subscribeThenComplete}, "A", Complete)
		subject.Next("B")
		subscribe(t, []Observable{subscribeThenComplete}, "A", "B", Complete)
		subject.Next("C")
		subscribe(t, []Observable{subscribeThenComplete}, "A", "B", "C", Complete)
		subject.Next("D")
		subscribe(t, []Observable{subscribeThenComplete}, "B", "C", "D", Complete)
		subject.Error(errTest)
		subscribe(t, []Observable{subscribeThenComplete}, errTest)
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
		subscribe(t, []Observable{subscribeThenComplete}, "A", Complete)
		time.Sleep(step(2))
		subject.Next("B")
		subscribe(t, []Observable{subscribeThenComplete}, "A", "B", Complete)
		time.Sleep(step(2))
		subject.Next("C")
		subscribe(t, []Observable{subscribeThenComplete}, "A", "B", "C", Complete)
		time.Sleep(step(2))
		subject.Next("D")
		subject.Complete()
		subscribe(t, []Observable{subscribeThenComplete}, "B", "C", "D", Complete)
		time.Sleep(step(2))
		subscribe(t, []Observable{subscribeThenComplete}, "C", "D", Complete)
		time.Sleep(step(2))
		subscribe(t, []Observable{subscribeThenComplete}, "D", Complete)
		time.Sleep(step(2))
		subscribe(t, []Observable{subscribeThenComplete}, Complete)
	})
}
