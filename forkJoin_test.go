package rx_test

import (
	"context"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestForkJoin(t *testing.T) {
	NewTestSuite(t).Case(
		rx.ForkJoin(
			rx.Just("A", "B", "C").Pipe(AddLatencyToValues(0, 3)),
			rx.Range(1, 5).Pipe(AddLatencyToValues(1, 2)),
			rx.Range(5, 9).Pipe(AddLatencyToValues(3, 1)),
		).Pipe(
			ToString(),
		),
		"[C 4 8]", Completed,
	).Case(
		rx.ForkJoin(
			rx.Just("A", "B", "C").Pipe(AddLatencyToValues(0, 3)),
			rx.Range(1, 5).Pipe(AddLatencyToValues(1, 2)),
			rx.Range(5, 9).Pipe(AddLatencyToValues(3, 1)),
			rx.Empty().Pipe(DelaySubscription(5)),
		).Pipe(
			ToString(),
		),
		Completed,
	).Case(
		rx.ForkJoin(
			rx.Just("A", "B", "C").Pipe(AddLatencyToValues(0, 3)),
			rx.Range(1, 5).Pipe(AddLatencyToValues(1, 2)),
			rx.Range(5, 9).Pipe(AddLatencyToValues(3, 1)),
			rx.Throw(ErrTest).Pipe(DelaySubscription(5)),
		).Pipe(
			ToString(),
		),
		ErrTest,
	).Case(
		rx.ForkJoin(),
		Completed,
	).TestAll()

	ctx, cancel := context.WithTimeout(context.Background(), Step(1))
	defer cancel()

	rx.ForkJoin(
		rx.Timer(Step(2)),
	).Subscribe(ctx, func(rx.Notification) { t.Fatal("should not happen") })
}
