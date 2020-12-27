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
		d := rx.UnicastLatest()

		rx.Just("A", "B", "C").Pipe(
			AddLatencyToNotifications(1, 1),
		).Subscribe(context.Background(), d.Observer)

		NewTestSuite(t).Case(
			d.Observable,
			"A", "B", "C", Completed,
		).Case(
			d.Observable,
			Completed,
		).TestAll()
	})

	t.Run("EarlyUnsubscribed", func(t *testing.T) {
		d := rx.UnicastLatest()

		rx.Just("A", "B", "C", "D").Pipe(
			AddLatencyToNotifications(1, 2),
		).Subscribe(context.Background(), d.Observer)

		NewTestSuite(t).Case(
			d.Observable.Pipe(operators.Take(2)),
			"A", "B", Completed,
		).Case(
			rx.Merge(d.Observable, d.Observable).Pipe(
				DelaySubscription(3),
			),
			rx.ErrDropped,
		).Case(
			d.Observable,
			"D", Completed,
		).Case(
			d.Observable,
			Completed,
		).TestAll()
	})

	t.Run("AfterComplete", func(t *testing.T) {
		d := rx.UnicastLatest()

		d.Complete()

		Test(t, d.Observable, Completed)

		d.Error(ErrTest)

		Test(t, d.Observable, Completed)
	})

	t.Run("AfterError", func(t *testing.T) {
		d := rx.UnicastLatest()

		d.Error(ErrTest)

		Test(t, d.Observable, ErrTest)

		d.Complete()

		Test(t, d.Observable, ErrTest)
	})

	t.Run("NilError", func(t *testing.T) {
		d := rx.UnicastLatest()

		rx.Throw(nil).Subscribe(context.Background(), d.Observer)

		Test(t, d.Observable, nil)
	})
}
