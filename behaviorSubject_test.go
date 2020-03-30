package rx_test

import (
	"context"
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestBehaviorSubject(t *testing.T) {
	sum := func(acc, val interface{}, idx int) interface{} {
		return acc.(int) + val.(int)
	}

	t.Run("Complete", func(t *testing.T) {
		subject := rx.NewBehaviorSubject(0)

		rx.Just(3, 4, 5).Pipe(
			AddLatencyToValues(1, 1),
		).Subscribe(context.Background(), subject.Observer)

		Subscribe(
			t,
			rx.Zip(
				subject.Observable,
				subject.Pipe(operators.Scan(sum)),
			).Pipe(ToString()),
			"[0 0]", "[3 3]", "[4 7]", "[5 12]", rx.Complete,
		)

		Subscribe(t, subject.Observable, 5, rx.Complete)
	})

	t.Run("Error", func(t *testing.T) {
		subject := rx.NewBehaviorSubject(0)

		rx.Concat(rx.Just(3, 4, 5), rx.Throw(ErrTest)).Pipe(
			AddLatencyToNotifications(1, 1),
		).Subscribe(context.Background(), subject.Observer)

		Subscribe(
			t,
			rx.Zip(
				subject.Observable,
				subject.Pipe(operators.Scan(sum)),
			).Pipe(ToString()),
			"[0 0]", "[3 3]", "[4 7]", "[5 12]", ErrTest,
		)

		Subscribe(t, subject.Observable, ErrTest)
	})
}
