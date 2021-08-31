package rx_test

import (
	"context"
	"runtime"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestMulticast(t *testing.T) {
	sum := func(acc, val interface{}, idx int) interface{} {
		return acc.(int) + val.(int)
	}

	t.Run("Normal", func(t *testing.T) {
		m := rx.Multicast()

		rx.Just(3, 4, 5).Pipe(
			AddLatencyToValues(1, 1),
		).Subscribe(context.Background(), m.Observer)

		NewTestSuite(t).Case(
			rx.Zip(
				m.Observable,
				m.Pipe(operators.Scan(sum)),
			).Pipe(
				ToString(),
			),
			"[3 3]", "[4 7]", "[5 12]", Completed,
		).Case(
			m.Observable,
			Completed,
		).TestAll()
	})

	t.Run("Error", func(t *testing.T) {
		m := rx.Multicast()

		rx.Concat(
			rx.Just(3, 4, 5),
			rx.Throw(ErrTest),
		).Pipe(
			AddLatencyToNotifications(1, 1),
		).Subscribe(context.Background(), m.Observer)

		NewTestSuite(t).Case(
			rx.Zip(
				m.Observable,
				m.Pipe(operators.Scan(sum)),
			).Pipe(
				ToString(),
			),
			"[3 3]", "[4 7]", "[5 12]", ErrTest,
		).Case(
			m.Observable,
			ErrTest,
		).TestAll()
	})

	t.Run("AfterComplete", func(t *testing.T) {
		m := rx.Multicast()

		m.Complete()

		Test(t, m.Observable, Completed)

		m.Error(ErrTest)

		Test(t, m.Observable, Completed)
	})

	t.Run("AfterError", func(t *testing.T) {
		m := rx.Multicast()

		m.Error(ErrTest)

		Test(t, m.Observable, ErrTest)

		m.Complete()

		Test(t, m.Observable, ErrTest)
	})

	t.Run("NilError", func(t *testing.T) {
		m := rx.Multicast()

		rx.Throw(nil).Subscribe(context.Background(), m.Observer)

		Test(t, m.Observable, nil)
	})

	t.Run("Finalizer", func(t *testing.T) {
		m := rx.Multicast()

		for i := 0; i < 10; i++ {
			m.Subscribe(context.Background(), rx.Noop)
		}

		runtime.GC()
	})
}
