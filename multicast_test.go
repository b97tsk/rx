package rx_test

import (
	"runtime"
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestMulticast(t *testing.T) {
	t.Parallel()

	t.Run("Buffer", func(t *testing.T) {
		t.Parallel()

		m := rx.MulticastBuffer[string](3)

		subscribeThenComplete := rx.NewObservable(
			func(c rx.Context, sink rx.Observer[string]) {
				c, sink = rx.Serialize(c, sink)
				m.Subscribe(c, sink)
				sink.Complete()
			},
		)

		m.Next("A")

		NewTestSuite[string](t).Case(subscribeThenComplete, "A", ErrComplete)

		m.Next("B")

		NewTestSuite[string](t).Case(subscribeThenComplete, "A", "B", ErrComplete)

		m.Next("C")

		NewTestSuite[string](t).Case(subscribeThenComplete, "A", "B", "C", ErrComplete)

		ctx, cancel := rx.NewBackgroundContext().WithTimeout(Step(2))
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

		NewTestSuite[string](t).Case(subscribeThenComplete, "B", "C", "D", ErrTest)
	})

	t.Run("BufferAll", func(t *testing.T) {
		t.Parallel()

		m := rx.MulticastBufferAll[string]()

		for _, v := range []string{"A", "B", "C"} {
			m.Next(v)
		}

		m.Complete()

		NewTestSuite[string](t).Case(
			m.Observable,
			"A", "B", "C", ErrComplete,
		).Case(
			rx.Pipe1(
				m.Observable,
				rx.OnNext(func(string) { panic(ErrTest) }),
			),
			rx.ErrOops, ErrTest,
		)
	})

	t.Run("Multicast", func(t *testing.T) {
		t.Parallel()

		m := rx.Multicast[string]()

		for _, v := range []string{"A", "B", "C"} {
			m.Next(v)
		}

		m.Complete()

		NewTestSuite[string](t).Case(m.Observable, ErrComplete)
	})

	t.Run("AfterComplete", func(t *testing.T) {
		t.Parallel()

		m := rx.Multicast[string]()

		m.Complete()

		NewTestSuite[string](t).Case(m.Observable, ErrComplete)

		m.Error(ErrTest)

		NewTestSuite[string](t).Case(m.Observable, ErrComplete)
	})

	t.Run("AfterError", func(t *testing.T) {
		t.Parallel()

		m := rx.Multicast[string]()

		m.Error(ErrTest)

		NewTestSuite[string](t).Case(m.Observable, ErrTest)

		m.Complete()

		NewTestSuite[string](t).Case(m.Observable, ErrTest)
	})

	t.Run("Finalizer", func(t *testing.T) {
		t.Parallel()

		c := make(chan struct{})

		m := rx.Multicast[string]()

		m.Subscribe(rx.NewBackgroundContext(), func(n rx.Notification[string]) {
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
