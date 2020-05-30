package rx_test

import (
	"context"
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestSubject(t *testing.T) {
	sum := func(acc, val interface{}, idx int) interface{} {
		return acc.(int) + val.(int)
	}

	t.Run("Completed", func(t *testing.T) {
		subject := rx.NewSubject()

		rx.Just(3, 4, 5).Pipe(
			AddLatencyToValues(1, 1),
		).Subscribe(context.Background(), subject.Observer)

		Subscribe(
			t,
			rx.Zip(
				subject.Observable,
				subject.Pipe(operators.Scan(sum)),
			).Pipe(ToString()),
			"[3 3]", "[4 7]", "[5 12]", rx.Completed,
		)

		Subscribe(t, subject.Observable, rx.Completed)
	})

	t.Run("Error", func(t *testing.T) {
		subject := rx.NewSubject()

		rx.Concat(rx.Just(3, 4, 5), rx.Throw(ErrTest)).Pipe(
			AddLatencyToNotifications(1, 1),
		).Subscribe(context.Background(), subject.Observer)

		Subscribe(
			t,
			rx.Zip(
				subject.Observable,
				subject.Pipe(operators.Scan(sum)),
			).Pipe(ToString()),
			"[3 3]", "[4 7]", "[5 12]", ErrTest,
		)

		Subscribe(t, subject.Observable, ErrTest)
	})
}
