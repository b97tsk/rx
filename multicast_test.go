package rx_test

import (
	"context"
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
		d := rx.Multicast()

		rx.Just(3, 4, 5).Pipe(
			AddLatencyToValues(1, 1),
		).Subscribe(context.Background(), d.Observer)

		NewTestSuite(t).Case(
			rx.Zip(
				d.Observable,
				d.Pipe(operators.Scan(sum)),
			).Pipe(
				ToString(),
			),
			"[3 3]", "[4 7]", "[5 12]", Completed,
		).Case(
			d.Observable,
			Completed,
		).TestAll()
	})

	t.Run("Error", func(t *testing.T) {
		d := rx.Multicast()

		rx.Concat(
			rx.Just(3, 4, 5),
			rx.Throw(ErrTest),
		).Pipe(
			AddLatencyToNotifications(1, 1),
		).Subscribe(context.Background(), d.Observer)

		NewTestSuite(t).Case(
			rx.Zip(
				d.Observable,
				d.Pipe(operators.Scan(sum)),
			).Pipe(
				ToString(),
			),
			"[3 3]", "[4 7]", "[5 12]", ErrTest,
		).Case(
			d.Observable,
			ErrTest,
		).TestAll()
	})

	t.Run("AfterComplete", func(t *testing.T) {
		d := rx.Multicast()

		d.Complete()

		Test(t, d.Observable, Completed)

		d.Error(ErrTest)

		Test(t, d.Observable, Completed)
	})

	t.Run("AfterError", func(t *testing.T) {
		d := rx.Multicast()

		d.Error(ErrTest)

		Test(t, d.Observable, ErrTest)

		d.Complete()

		Test(t, d.Observable, ErrTest)
	})

	t.Run("NilError", func(t *testing.T) {
		d := rx.Multicast()

		rx.Throw(nil).Subscribe(context.Background(), d.Observer)

		Test(t, d.Observable, nil)
	})
}
