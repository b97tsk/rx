package rx_test

import (
	"context"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestCombineLatest(t *testing.T) {
	NewTestSuite(t).Case(
		rx.CombineLatest(
			rx.Just("A", "B").Pipe(AddLatencyToValues(1, 3)),
			rx.Just("C", "D").Pipe(AddLatencyToValues(2, 4)),
			rx.Just("E", "F").Pipe(AddLatencyToValues(3, 5)),
		).Pipe(
			ToString(),
		),
		"[A C E]", "[B C E]", "[B D E]", "[B D F]", Completed,
	).Case(
		rx.CombineLatest(
			rx.Just("A", "B").Pipe(AddLatencyToValues(1, 3)),
			rx.Just("C", "D").Pipe(AddLatencyToValues(2, 4)),
			rx.Throw(ErrTest).Pipe(DelaySubscription(5)),
		).Pipe(
			ToString(),
		),
		ErrTest,
	).Case(
		rx.CombineLatest(),
		Completed,
	).TestAll()

	ctx, cancel := context.WithTimeout(context.Background(), Step(1))
	defer cancel()

	rx.CombineLatest(
		rx.Timer(Step(2)),
	).Subscribe(ctx, func(rx.Notification) { t.Fatal("should not happen") })
}
