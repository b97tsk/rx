package rx_test

import (
	"context"
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestMulticastReplay(t *testing.T) {
	t.Run("BufferSize", func(t *testing.T) {
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

		ctx, cancel := context.WithTimeout(context.Background(), Step(2))
		defer cancel()

		go d.Observable.Pipe(
			operators.DoOnNext(
				func(interface{}) {
					time.Sleep(Step(2))
				},
			),
		).Subscribe(ctx, rx.Noop)

		time.Sleep(Step(1))

		d.Next("D")

		time.Sleep(Step(2))

		Test(t, subscribeThenComplete, "B", "C", "D", Completed)

		d.Error(ErrTest)

		Test(t, subscribeThenComplete, ErrTest)
	})

	t.Run("WindowTime", func(t *testing.T) {
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
	})

	t.Run("AfterComplete", func(t *testing.T) {
		d := rx.MulticastReplay(nil)

		d.Complete()

		Test(t, d.Observable, Completed)

		d.Error(ErrTest)

		Test(t, d.Observable, Completed)
	})

	t.Run("AfterError", func(t *testing.T) {
		d := rx.MulticastReplay(nil)

		d.Error(ErrTest)

		Test(t, d.Observable, ErrTest)

		d.Complete()

		Test(t, d.Observable, ErrTest)
	})

	t.Run("NilError", func(t *testing.T) {
		d := rx.MulticastReplayFactory(nil)()

		rx.Throw(nil).Subscribe(context.Background(), d.Observer)

		Test(t, d.Observable, nil)
	})
}
