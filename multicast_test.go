package rx_test

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestMulticast(t *testing.T) {
	t.Parallel()

	sum := func(v1, v2 int) int {
		return v1 + v2
	}

	toString := func(v1, v2 int) string {
		return fmt.Sprintf("[%v %v]", v1, v2)
	}

	t.Run("Normal", func(t *testing.T) {
		t.Parallel()

		m := rx.Multicast[int]()

		rx.Pipe1(
			rx.Just(3, 4, 5),
			AddLatencyToValues[int](1, 1),
		).Subscribe(context.Background(), m.Observer)

		NewTestSuite[string](t).Case(
			rx.Zip2(
				m.Observable,
				rx.Pipe1(m.Observable, rx.Scan(0, sum)),
				toString,
			),
			"[3 3]", "[4 7]", "[5 12]", ErrCompleted,
		).Case(
			rx.Pipe1(
				m.Observable,
				ToString[int](),
			),
			ErrCompleted,
		)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		m := rx.Multicast[int]()

		rx.Pipe1(
			rx.Concat(
				rx.Just(3, 4, 5),
				rx.Throw[int](ErrTest),
			),
			AddLatencyToNotifications[int](1, 1),
		).Subscribe(context.Background(), m.Observer)

		NewTestSuite[string](t).Case(
			rx.Zip2(
				m.Observable,
				rx.Pipe1(m.Observable, rx.Scan(0, sum)),
				toString,
			),
			"[3 3]", "[4 7]", "[5 12]", ErrTest,
		).Case(
			rx.Pipe1(
				m.Observable,
				ToString[int](),
			),
			ErrTest,
		)
	})

	t.Run("AfterComplete", func(t *testing.T) {
		t.Parallel()

		m := rx.Multicast[int]()

		m.Complete()

		NewTestSuite[int](t).Case(m.Observable, ErrCompleted)

		m.Error(ErrTest)

		NewTestSuite[int](t).Case(m.Observable, ErrCompleted)
	})

	t.Run("AfterError", func(t *testing.T) {
		t.Parallel()

		m := rx.Multicast[int]()

		m.Error(ErrTest)

		NewTestSuite[int](t).Case(m.Observable, ErrTest)

		m.Complete()

		NewTestSuite[int](t).Case(m.Observable, ErrTest)
	})

	t.Run("NilError", func(t *testing.T) {
		t.Parallel()

		m := rx.Multicast[int]()

		m.Error(nil)

		NewTestSuite[int](t).Case(m.Observable, nil)
	})

	t.Run("Finalizer", func(t *testing.T) {
		t.Parallel()

		c := make(chan struct{})

		m := rx.Multicast[int]()
		m.Subscribe(context.Background(), func(n rx.Notification[int]) {
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
