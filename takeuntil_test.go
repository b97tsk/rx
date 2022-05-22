package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestTakeUntil(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](0, 2),
			rx.TakeUntil[string](rx.Just(42)),
		),
		ErrCompleted,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](0, 2),
			rx.TakeUntil[string](rx.Empty[int]()),
		),
		"A", "B", "C", ErrCompleted,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](0, 2),
			rx.TakeUntil[string](rx.Never[int]()),
		),
		"A", "B", "C", ErrCompleted,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](0, 2),
			rx.TakeUntil[string](rx.Throw[int](ErrTest)),
		),
		ErrTest,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](0, 2),
			rx.TakeUntil[string](
				rx.Pipe(
					rx.Just(42),
					DelaySubscription[int](3),
				),
			),
		),
		"A", "B", ErrCompleted,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](0, 2),
			rx.TakeUntil[string](
				rx.Pipe(
					rx.Empty[int](),
					DelaySubscription[int](3),
				),
			),
		),
		"A", "B", "C", ErrCompleted,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](0, 2),
			rx.TakeUntil[string](
				rx.Pipe(
					rx.Throw[int](ErrTest),
					DelaySubscription[int](3),
				),
			),
		),
		"A", "B", ErrTest,
	)
}
