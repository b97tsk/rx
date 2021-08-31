package rx_test

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestMulticastReplay(t *testing.T) {
	t.Run("BufferSize", func(t *testing.T) {
		m := rx.MulticastReplay(&rx.ReplayOptions{BufferSize: 3})

		subscribeThenComplete := rx.Observable(
			func(ctx context.Context, sink rx.Observer) {
				sink = sink.Mutex()
				m.Subscribe(ctx, sink)
				sink.Complete()
			},
		)

		m.Next("A")

		Test(t, subscribeThenComplete, "A", Completed)

		m.Next("B")

		Test(t, subscribeThenComplete, "A", "B", Completed)

		m.Next("C")

		Test(t, subscribeThenComplete, "A", "B", "C", Completed)

		ctx, cancel := context.WithTimeout(context.Background(), Step(2))
		defer cancel()

		go m.Observable.Pipe(
			operators.DoOnNext(
				func(interface{}) {
					time.Sleep(Step(2))
				},
			),
		).Subscribe(ctx, rx.Noop)

		time.Sleep(Step(1))

		m.Next("D")

		time.Sleep(Step(2))

		Test(t, subscribeThenComplete, "B", "C", "D", Completed)

		m.Error(ErrTest)

		Test(t, subscribeThenComplete, ErrTest)
	})

	t.Run("WindowTime", func(t *testing.T) {
		m := rx.MulticastReplay(&rx.ReplayOptions{WindowTime: Step(5)})

		subscribeThenComplete := rx.Observable(
			func(ctx context.Context, sink rx.Observer) {
				sink = sink.Mutex()
				m.Subscribe(ctx, sink)
				sink.Complete()
			},
		)

		m.Next("A")

		Test(t, subscribeThenComplete, "A", Completed)

		time.Sleep(Step(2))
		m.Next("B")

		Test(t, subscribeThenComplete, "A", "B", Completed)

		time.Sleep(Step(2))
		m.Next("C")

		Test(t, subscribeThenComplete, "A", "B", "C", Completed)

		time.Sleep(Step(2))
		m.Next("D")
		m.Complete()

		Test(t, subscribeThenComplete, "B", "C", "D", Completed)

		time.Sleep(Step(2))

		Test(t, subscribeThenComplete, "C", "D", Completed)

		time.Sleep(Step(2))

		Test(t, subscribeThenComplete, "D", Completed)

		time.Sleep(Step(2))

		Test(t, subscribeThenComplete, Completed)
	})

	t.Run("AfterComplete", func(t *testing.T) {
		m := rx.MulticastReplay(nil)

		m.Complete()

		Test(t, m.Observable, Completed)

		m.Error(ErrTest)

		Test(t, m.Observable, Completed)
	})

	t.Run("AfterError", func(t *testing.T) {
		m := rx.MulticastReplay(nil)

		m.Error(ErrTest)

		Test(t, m.Observable, ErrTest)

		m.Complete()

		Test(t, m.Observable, ErrTest)
	})

	t.Run("NilError", func(t *testing.T) {
		m := rx.MulticastReplayFactory(nil)()

		rx.Throw(nil).Subscribe(context.Background(), m.Observer)

		Test(t, m.Observable, nil)
	})

	t.Run("Finalizer", func(t *testing.T) {
		m := rx.MulticastReplay(nil)

		for i := 0; i < 10; i++ {
			m.Subscribe(context.Background(), rx.Noop)
		}

		runtime.GC()
	})
}
