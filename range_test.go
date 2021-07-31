package rx_test

import (
	"context"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestRange(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Range(1, 4),
		1, 2, 3, Completed,
	).Case(
		rx.Range(4, 1),
		Completed,
	).TestAll()

	ctx, cancel := context.WithTimeout(context.Background(), Step(1))
	defer cancel()

	rx.Range(1, 4).Pipe(
		AddLatencyToValues(2, 2),
	).BlockingSubscribe(ctx, func(rx.Notification) { t.Fatal("should not happen") })
}

func TestIota(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Iota(1).Pipe(
			operators.Take(3),
		),
		1, 2, 3, Completed,
	).TestAll()

	ctx, cancel := context.WithTimeout(context.Background(), Step(1))
	defer cancel()

	rx.Iota(1).Pipe(
		operators.Take(3),
		AddLatencyToValues(2, 2),
	).BlockingSubscribe(ctx, func(rx.Notification) { t.Fatal("should not happen") })
}
