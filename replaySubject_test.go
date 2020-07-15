package rx_test

import (
	"context"
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestReplaySubject(t *testing.T) {
	t.Run("#1", func(t *testing.T) {
		subject := rx.NewReplaySubject(0)
		subject.SetBufferSize(3)
		subscribeThenComplete := rx.Observable(
			func(ctx context.Context, sink rx.Observer) {
				sink = sink.Mutex()
				subject.Subscribe(ctx, sink)
				sink.Complete()
			},
		)
		subject.Next("A")
		Subscribe(t, subscribeThenComplete, "A", Completed)
		subject.Next("B")
		Subscribe(t, subscribeThenComplete, "A", "B", Completed)
		subject.Next("C")
		Subscribe(t, subscribeThenComplete, "A", "B", "C", Completed)
		subject.Next("D")
		Subscribe(t, subscribeThenComplete, "B", "C", "D", Completed)
		subject.Error(ErrTest)
		Subscribe(t, subscribeThenComplete, ErrTest)
	})
	t.Run("#2", func(t *testing.T) {
		subject := rx.NewReplaySubject(0)
		subject.SetWindowTime(Step(5))
		subscribeThenComplete := rx.Observable(
			func(ctx context.Context, sink rx.Observer) {
				sink = sink.Mutex()
				subject.Subscribe(ctx, sink)
				sink.Complete()
			},
		)
		subject.Next("A")
		Subscribe(t, subscribeThenComplete, "A", Completed)
		time.Sleep(Step(2))
		subject.Next("B")
		Subscribe(t, subscribeThenComplete, "A", "B", Completed)
		time.Sleep(Step(2))
		subject.Next("C")
		Subscribe(t, subscribeThenComplete, "A", "B", "C", Completed)
		time.Sleep(Step(2))
		subject.Next("D")
		subject.Complete()
		Subscribe(t, subscribeThenComplete, "B", "C", "D", Completed)
		time.Sleep(Step(2))
		Subscribe(t, subscribeThenComplete, "C", "D", Completed)
		time.Sleep(Step(2))
		Subscribe(t, subscribeThenComplete, "D", Completed)
		time.Sleep(Step(2))
		Subscribe(t, subscribeThenComplete, Completed)
	})
}
