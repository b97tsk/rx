package subject_test

import (
	"context"
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/subject"
)

func TestReplaySubject(t *testing.T) {
	t.Run("#1", func(t *testing.T) {
		subject1 := subject.NewReplaySubject(0)
		subject1.SetBufferSize(3)
		subscribeThenComplete := rx.Observable(
			func(ctx context.Context, sink rx.Observer) {
				sink = sink.Mutex()
				subject1.Subscribe(ctx, sink)
				sink.Complete()
			},
		)
		subject1.Next("A")
		Subscribe(t, subscribeThenComplete, "A", Completed)
		subject1.Next("B")
		Subscribe(t, subscribeThenComplete, "A", "B", Completed)
		subject1.Next("C")
		Subscribe(t, subscribeThenComplete, "A", "B", "C", Completed)
		subject1.Next("D")
		Subscribe(t, subscribeThenComplete, "B", "C", "D", Completed)
		subject1.Error(ErrTest)
		Subscribe(t, subscribeThenComplete, ErrTest)
	})
	t.Run("#2", func(t *testing.T) {
		subject1 := subject.NewReplaySubject(0)
		subject1.SetWindowTime(Step(5))
		subscribeThenComplete := rx.Observable(
			func(ctx context.Context, sink rx.Observer) {
				sink = sink.Mutex()
				subject1.Subscribe(ctx, sink)
				sink.Complete()
			},
		)
		subject1.Next("A")
		Subscribe(t, subscribeThenComplete, "A", Completed)
		time.Sleep(Step(2))
		subject1.Next("B")
		Subscribe(t, subscribeThenComplete, "A", "B", Completed)
		time.Sleep(Step(2))
		subject1.Next("C")
		Subscribe(t, subscribeThenComplete, "A", "B", "C", Completed)
		time.Sleep(Step(2))
		subject1.Next("D")
		subject1.Complete()
		Subscribe(t, subscribeThenComplete, "B", "C", "D", Completed)
		time.Sleep(Step(2))
		Subscribe(t, subscribeThenComplete, "C", "D", Completed)
		time.Sleep(Step(2))
		Subscribe(t, subscribeThenComplete, "D", Completed)
		time.Sleep(Step(2))
		Subscribe(t, subscribeThenComplete, Completed)
	})
}
