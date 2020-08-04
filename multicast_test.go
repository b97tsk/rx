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

	t.Run("Completed", func(t *testing.T) {
		d := rx.Multicast()

		rx.Just(3, 4, 5).Pipe(
			AddLatencyToValues(1, 1),
		).Subscribe(context.Background(), d.Observer)

		Subscribe(
			t,
			rx.Zip(
				d.Observable,
				d.Pipe(operators.Scan(sum)),
			).Pipe(ToString()),
			"[3 3]", "[4 7]", "[5 12]", Completed,
		)

		Subscribe(t, d.Observable, Completed)
	})

	t.Run("Error", func(t *testing.T) {
		d := rx.Multicast()

		rx.Concat(rx.Just(3, 4, 5), rx.Throw(ErrTest)).Pipe(
			AddLatencyToNotifications(1, 1),
		).Subscribe(context.Background(), d.Observer)

		Subscribe(
			t,
			rx.Zip(
				d.Observable,
				d.Pipe(operators.Scan(sum)),
			).Pipe(ToString()),
			"[3 3]", "[4 7]", "[5 12]", ErrTest,
		)

		Subscribe(t, d.Observable, ErrTest)
	})
}
