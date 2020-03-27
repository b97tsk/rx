package rx_test

import (
	"context"
	"testing"

	. "github.com/b97tsk/rx"
)

func TestBehaviorSubject(t *testing.T) {
	sum := func(acc, val interface{}, idx int) interface{} {
		return acc.(int) + val.(int)
	}

	t.Run("Complete", func(t *testing.T) {
		subject := NewBehaviorSubject(0)

		Just(3, 4, 5).Pipe(
			addLatencyToValue(1, 1),
		).Subscribe(context.Background(), subject.Observer)

		subscribe(
			t,
			Zip(
				subject.Observable,
				subject.Pipe(operators.Scan(sum)),
			).Pipe(toString),
			"[0 0]", "[3 3]", "[4 7]", "[5 12]", Complete,
		)

		subscribe(t, subject.Observable, 5, Complete)
	})

	t.Run("Error", func(t *testing.T) {
		subject := NewBehaviorSubject(0)

		Concat(Just(3, 4, 5), Throw(errTest)).Pipe(
			addLatencyToNotification(1, 1),
		).Subscribe(context.Background(), subject.Observer)

		subscribe(
			t,
			Zip(
				subject.Observable,
				subject.Pipe(operators.Scan(sum)),
			).Pipe(toString),
			"[0 0]", "[3 3]", "[4 7]", "[5 12]", errTest,
		)

		subscribe(t, subject.Observable, errTest)
	})
}
