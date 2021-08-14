package rx_test

import (
	"context"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestUnicastLatest(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		u := rx.UnicastLatest()

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
		u := rx.UnicastLatest()

		rx.Just("A", "B", "C", "D").Pipe(
			AddLatencyToNotifications(1, 2),
		).Subscribe(context.Background(), u.Observer)

		NewTestSuite(t).Case(
			u.Observable.Pipe(operators.Take(2)),
			"A", "B", Completed,
		).Case(
			rx.Merge(u.Observable, u.Observable).Pipe(
				DelaySubscription(3),
			),
			rx.ErrDropped,
		).Case(
			u.Observable,
			"D", Completed,
		).Case(
			u.Observable,
			Completed,
		).TestAll()
	})

	t.Run("AfterComplete", func(t *testing.T) {
		u := rx.UnicastLatest()

		u.Complete()

		Test(t, u.Observable, Completed)

		u.Error(ErrTest)

		Test(t, u.Observable, Completed)
	})

	t.Run("AfterError", func(t *testing.T) {
		u := rx.UnicastLatest()

		u.Error(ErrTest)

		Test(t, u.Observable, ErrTest)

		u.Complete()

		Test(t, u.Observable, ErrTest)
	})

	t.Run("NilError", func(t *testing.T) {
		u := rx.UnicastLatest()

		rx.Throw(nil).Subscribe(context.Background(), u.Observer)

		Test(t, u.Observable, nil)
	})
}
