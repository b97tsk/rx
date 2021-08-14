package rx_test

import (
	"context"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestUnicast(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		u := rx.Unicast()

		rx.Just("A", "B", "C").Pipe(
			AddLatencyToNotifications(1, 1),
		).Subscribe(context.Background(), u.Observer)

		NewTestSuite(t).Case(
			u.Observable,
			"A", "B", "C", Completed,
		).Case(
			u.Observable,
			Completed,
		).TestAll()
	})

	t.Run("EarlyUnsubscribed", func(t *testing.T) {
		u := rx.Unicast()

		rx.Just("A", "B", "C").Pipe(
			AddLatencyToNotifications(1, 1),
		).Subscribe(context.Background(), u.Observer)

		NewTestSuite(t).Case(
			u.Observable.Pipe(operators.Take(2)),
			"A", "B", Completed,
		).Case(
			u.Observable,
			rx.ErrDropped,
		).Case(
			u.Observable.Pipe(DelaySubscription(3)),
			Completed,
		).Case(
			u.Observable,
			Completed,
		).TestAll()
	})

	t.Run("AfterComplete", func(t *testing.T) {
		u := rx.Unicast()

		u.Complete()

		Test(t, u.Observable, Completed)

		u.Error(ErrTest)

		Test(t, u.Observable, Completed)
	})

	t.Run("AfterError", func(t *testing.T) {
		u := rx.Unicast()

		u.Error(ErrTest)

		Test(t, u.Observable, ErrTest)

		u.Complete()

		Test(t, u.Observable, ErrTest)
	})

	t.Run("NilError", func(t *testing.T) {
		u := rx.Unicast()

		rx.Throw(nil).Subscribe(context.Background(), u.Observer)

		Test(t, u.Observable, nil)
	})
}
