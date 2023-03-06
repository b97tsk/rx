package rx_test

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestMulticastReplay(t *testing.T) {
	t.Parallel()

	t.Run("BufferSize", func(t *testing.T) {
		t.Parallel()

		m := rx.MulticastReplay[string](&rx.ReplayConfig{BufferSize: 3})

		subscribeThenComplete := rx.NewObservable(
			func(ctx context.Context, sink rx.Observer[string]) {
				ctx, cancel := context.WithCancel(ctx)
				sink = sink.WithMutex()
				m.Subscribe(ctx, sink)
				sink.Complete()
				cancel()
			},
		)

		m.Next("A")

		NewTestSuite[string](t).Case(subscribeThenComplete, "A", ErrComplete)

		m.Next("B")

		NewTestSuite[string](t).Case(subscribeThenComplete, "A", "B", ErrComplete)

		m.Next("C")

		NewTestSuite[string](t).Case(subscribeThenComplete, "A", "B", "C", ErrComplete)

		ctx, cancel := context.WithTimeout(context.Background(), Step(2))
		defer cancel()

		go rx.Pipe1(
			m.Observable,
			rx.OnNext(
				func(string) {
					time.Sleep(Step(2))
				},
			),
		).Subscribe(ctx, rx.Noop[string])

		time.Sleep(Step(1))

		m.Next("D")

		time.Sleep(Step(2))

		NewTestSuite[string](t).Case(subscribeThenComplete, "B", "C", "D", ErrComplete)

		m.Error(ErrTest)

		NewTestSuite[string](t).Case(subscribeThenComplete, ErrTest)
	})

	t.Run("WindowTime", func(t *testing.T) {
		t.Parallel()

		m := rx.MulticastReplay[string](&rx.ReplayConfig{WindowTime: Step(5)})

		subscribeThenComplete := rx.NewObservable(
			func(ctx context.Context, sink rx.Observer[string]) {
				ctx, cancel := context.WithCancel(ctx)
				sink = sink.WithMutex()
				m.Subscribe(ctx, sink)
				sink.Complete()
				cancel()
			},
		)

		m.Next("A")

		NewTestSuite[string](t).Case(subscribeThenComplete, "A", ErrComplete)

		time.Sleep(Step(2))
		m.Next("B")

		NewTestSuite[string](t).Case(subscribeThenComplete, "A", "B", ErrComplete)

		time.Sleep(Step(2))
		m.Next("C")

		NewTestSuite[string](t).Case(subscribeThenComplete, "A", "B", "C", ErrComplete)

		time.Sleep(Step(2))
		m.Next("D")
		m.Complete()

		NewTestSuite[string](t).Case(subscribeThenComplete, "B", "C", "D", ErrComplete)

		time.Sleep(Step(2))

		NewTestSuite[string](t).Case(subscribeThenComplete, "C", "D", ErrComplete)

		time.Sleep(Step(2))

		NewTestSuite[string](t).Case(subscribeThenComplete, "D", ErrComplete)

		time.Sleep(Step(2))

		NewTestSuite[string](t).Case(subscribeThenComplete, ErrComplete)
	})

	t.Run("AfterComplete", func(t *testing.T) {
		t.Parallel()

		m := rx.MulticastReplay[string](nil)

		m.Complete()

		NewTestSuite[string](t).Case(m.Observable, ErrComplete)

		m.Error(ErrTest)

		NewTestSuite[string](t).Case(m.Observable, ErrComplete)
	})

	t.Run("AfterError", func(t *testing.T) {
		t.Parallel()

		m := rx.MulticastReplay[string](nil)

		m.Error(ErrTest)

		NewTestSuite[string](t).Case(m.Observable, ErrTest)

		m.Complete()

		NewTestSuite[string](t).Case(m.Observable, ErrTest)
	})

	t.Run("NilError", func(t *testing.T) {
		t.Parallel()

		m := rx.MulticastReplay[string](nil)

		m.Error(nil)

		NewTestSuite[string](t).Case(m.Observable, nil)
	})

	t.Run("Finalizer", func(t *testing.T) {
		t.Parallel()

		c := make(chan struct{})

		m := rx.MulticastReplay[string](nil)
		m.Subscribe(context.Background(), func(n rx.Notification[string]) {
			if n.Error != rx.ErrFinalized {
				panic("want rx.ErrFinalized, but got something else")
			}

			close(c)
		})

		runtime.GC()

		select {
		case <-c:
		case <-time.After(5 * time.Second):
			t.Fatal("timeout waiting for running finalizer")
		}
	})
}
