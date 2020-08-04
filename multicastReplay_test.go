package rx_test

import (
	"context"
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestMulticastReplay(t *testing.T) {
	t.Run("#1", func(t *testing.T) {
		d := rx.MulticastReplay(&rx.ReplayOptions{BufferSize: 3})
		subscribeThenComplete := rx.Observable(
			func(ctx context.Context, sink rx.Observer) {
				sink = sink.Mutex()
				d.Subscribe(ctx, sink)
				sink.Complete()
			},
		)
		d.Next("A")
		Subscribe(t, subscribeThenComplete, "A", Completed)
		d.Next("B")
		Subscribe(t, subscribeThenComplete, "A", "B", Completed)
		d.Next("C")
		Subscribe(t, subscribeThenComplete, "A", "B", "C", Completed)
		d.Next("D")
		Subscribe(t, subscribeThenComplete, "B", "C", "D", Completed)
		d.Error(ErrTest)
		Subscribe(t, subscribeThenComplete, ErrTest)
	})
	t.Run("#2", func(t *testing.T) {
		d := rx.MulticastReplay(&rx.ReplayOptions{WindowTime: Step(5)})
		subscribeThenComplete := rx.Observable(
			func(ctx context.Context, sink rx.Observer) {
				sink = sink.Mutex()
				d.Subscribe(ctx, sink)
				sink.Complete()
			},
		)
		d.Next("A")
		Subscribe(t, subscribeThenComplete, "A", Completed)
		time.Sleep(Step(2))
		d.Next("B")
		Subscribe(t, subscribeThenComplete, "A", "B", Completed)
		time.Sleep(Step(2))
		d.Next("C")
		Subscribe(t, subscribeThenComplete, "A", "B", "C", Completed)
		time.Sleep(Step(2))
		d.Next("D")
		d.Complete()
		Subscribe(t, subscribeThenComplete, "B", "C", "D", Completed)
		time.Sleep(Step(2))
		Subscribe(t, subscribeThenComplete, "C", "D", Completed)
		time.Sleep(Step(2))
		Subscribe(t, subscribeThenComplete, "D", Completed)
		time.Sleep(Step(2))
		Subscribe(t, subscribeThenComplete, Completed)
	})
}
