package rx_test

import (
	"context"
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/testing"
)

func TestReplaySubject(t *testing.T) {
	t.Run("#1", func(t *testing.T) {
		subject := rx.NewReplaySubject(3, 0)
		subscribeThenComplete := rx.Create(
			func(ctx context.Context, sink rx.Observer) {
				sink = rx.Mutex(sink)
				subject.Subscribe(ctx, sink)
				sink.Complete()
			},
		)
		subject.Next("A")
		Subscribe(t, subscribeThenComplete, "A", rx.Complete)
		subject.Next("B")
		Subscribe(t, subscribeThenComplete, "A", "B", rx.Complete)
		subject.Next("C")
		Subscribe(t, subscribeThenComplete, "A", "B", "C", rx.Complete)
		subject.Next("D")
		Subscribe(t, subscribeThenComplete, "B", "C", "D", rx.Complete)
		subject.Error(ErrTest)
		Subscribe(t, subscribeThenComplete, ErrTest)
	})
	t.Run("#2", func(t *testing.T) {
		subject := rx.NewReplaySubject(0, Step(5))
		subscribeThenComplete := rx.Create(
			func(ctx context.Context, sink rx.Observer) {
				sink = rx.Mutex(sink)
				subject.Subscribe(ctx, sink)
				sink.Complete()
			},
		)
		subject.Next("A")
		Subscribe(t, subscribeThenComplete, "A", rx.Complete)
		time.Sleep(Step(2))
		subject.Next("B")
		Subscribe(t, subscribeThenComplete, "A", "B", rx.Complete)
		time.Sleep(Step(2))
		subject.Next("C")
		Subscribe(t, subscribeThenComplete, "A", "B", "C", rx.Complete)
		time.Sleep(Step(2))
		subject.Next("D")
		subject.Complete()
		Subscribe(t, subscribeThenComplete, "B", "C", "D", rx.Complete)
		time.Sleep(Step(2))
		Subscribe(t, subscribeThenComplete, "C", "D", rx.Complete)
		time.Sleep(Step(2))
		Subscribe(t, subscribeThenComplete, "D", rx.Complete)
		time.Sleep(Step(2))
		Subscribe(t, subscribeThenComplete, rx.Complete)
	})
}
