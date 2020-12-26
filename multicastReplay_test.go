package rx_test

import (
	"context"
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestMulticastReplay1(t *testing.T) {
	d := rx.MulticastReplay(&rx.ReplayOptions{BufferSize: 3})

	subscribeThenComplete := rx.Observable(
		func(ctx context.Context, sink rx.Observer) {
			sink = sink.Mutex()
			d.Subscribe(ctx, sink)
			sink.Complete()
		},
	)

	d.Next("A")

	Test(t, subscribeThenComplete, "A", Completed)

	d.Next("B")

	Test(t, subscribeThenComplete, "A", "B", Completed)

	d.Next("C")

	Test(t, subscribeThenComplete, "A", "B", "C", Completed)

	d.Next("D")

	Test(t, subscribeThenComplete, "B", "C", "D", Completed)

	d.Error(ErrTest)

	Test(t, subscribeThenComplete, ErrTest)
}

func TestMulticastReplay2(t *testing.T) {
	d := rx.MulticastReplay(&rx.ReplayOptions{WindowTime: Step(5)})

	subscribeThenComplete := rx.Observable(
		func(ctx context.Context, sink rx.Observer) {
			sink = sink.Mutex()
			d.Subscribe(ctx, sink)
			sink.Complete()
		},
	)

	d.Next("A")

	Test(t, subscribeThenComplete, "A", Completed)

	time.Sleep(Step(2))
	d.Next("B")

	Test(t, subscribeThenComplete, "A", "B", Completed)

	time.Sleep(Step(2))
	d.Next("C")

	Test(t, subscribeThenComplete, "A", "B", "C", Completed)

	time.Sleep(Step(2))
	d.Next("D")
	d.Complete()

	Test(t, subscribeThenComplete, "B", "C", "D", Completed)

	time.Sleep(Step(2))

	Test(t, subscribeThenComplete, "C", "D", Completed)

	time.Sleep(Step(2))

	Test(t, subscribeThenComplete, "D", Completed)

	time.Sleep(Step(2))

	Test(t, subscribeThenComplete, Completed)
}
