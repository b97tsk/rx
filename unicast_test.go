package rx_test

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestUnicast(t *testing.T) {
	t.Parallel()

	t.Run("SubscribeEarly", func(t *testing.T) {
		t.Parallel()

		var s []string

		u := rx.UnicastBuffer[string](3)

		u.Subscribe(rx.NewBackgroundContext(), func(n rx.Notification[string]) {
			if n.Kind == rx.KindNext {
				s = append(s, n.Value)
			}
		})

		u.Next("A")
		u.Next("B")
		u.Next("C")
		u.Next("D")
		u.Next("E")
		u.Complete()

		NewTestSuite[string](t).Case(
			rx.FromSlice(s),
			"A", "B", "C", "D", "E", ErrComplete,
		)
	})

	t.Run("CompleteEarly", func(t *testing.T) {
		t.Parallel()

		u := rx.UnicastBuffer[string](3)

		u.Next("A")
		u.Next("B")
		u.Next("C")
		u.Next("D")
		u.Next("E")
		u.Complete()

		NewTestSuite[string](t).Case(
			u.Observable,
			"C", "D", "E", ErrComplete,
		).Case(
			u.Observable,
			rx.ErrUnicast,
		)
	})

	t.Run("CompleteEarly #2", func(t *testing.T) {
		t.Parallel()

		u := rx.UnicastBuffer[string](5)

		u.Next("A")
		u.Next("B")
		u.Next("C")
		u.Next("D")
		u.Next("E")
		u.Complete()

		NewTestSuite[string](t).Case(
			rx.Pipe1(
				u.Observable,
				rx.Take[string](3),
			),
			"A", "B", "C", ErrComplete,
		).Case(
			u.Observable,
			rx.ErrUnicast,
		)
	})

	t.Run("ContextCancel", func(t *testing.T) {
		t.Parallel()

		u := rx.UnicastBufferAll[string]()
		defer u.Complete()

		u.Next("A")
		u.Next("B")
		u.Next("C")

		ctx, cancel := rx.NewBackgroundContext().WithTimeout(Step(1))
		defer cancel()

		NewTestSuite[string](t).WithContext(ctx).Case(
			u.Observable,
			"A", "B", "C", context.DeadlineExceeded,
		).Case(
			u.Observable,
			rx.ErrUnicast,
		)
	})

	t.Run("Oops", func(t *testing.T) {
		t.Parallel()

		u := rx.UnicastBufferAll[string]()

		u.Next("A")
		u.Next("B")
		u.Next("C")
		u.Complete()

		NewTestSuite[string](t).Case(
			rx.Pipe1(
				u.Observable,
				rx.OnNext(func(string) { panic(ErrTest) }),
			),
			rx.ErrOops, ErrTest,
		).Case(
			u.Observable,
			rx.ErrUnicast,
		)
	})

	t.Run("Empty", func(t *testing.T) {
		t.Parallel()

		u := rx.Unicast[string]()

		u.Complete()
		u.Error(ErrTest) // No-op.

		NewTestSuite[string](t).Case(
			u.Observable,
			ErrComplete,
		).Case(
			u.Observable,
			rx.ErrUnicast,
		)
	})

	t.Run("Empty #2", func(t *testing.T) {
		t.Parallel()

		u := rx.Unicast[string]()

		u.Error(ErrTest)
		u.Complete() // No-op.

		NewTestSuite[string](t).Case(
			u.Observable,
			ErrTest,
		).Case(
			u.Observable,
			rx.ErrUnicast,
		)
	})

	t.Run("Finalizer", func(t *testing.T) {
		t.Parallel()

		c := make(chan struct{})

		u := rx.Unicast[string]()

		u.Subscribe(rx.NewBackgroundContext(), func(n rx.Notification[string]) {
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
