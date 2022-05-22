package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestSkipUntil(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](0, 2),
			rx.SkipUntil[string](rx.Just(42)),
		),
		"A", "B", "C", ErrCompleted,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](0, 2),
			rx.SkipUntil[string](rx.Empty[int]()),
		),
		ErrCompleted,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](0, 2),
			rx.SkipUntil[string](rx.Never[int]()),
		),
		ErrCompleted,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](0, 2),
			rx.SkipUntil[string](rx.Throw[int](ErrTest)),
		),
		ErrTest,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](0, 2),
			rx.SkipUntil[string](
				rx.Pipe(
					rx.Just(42),
					DelaySubscription[int](3),
				),
			),
		),
		"C", ErrCompleted,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](0, 2),
			rx.SkipUntil[string](
				rx.Pipe(
					rx.Empty[int](),
					DelaySubscription[int](3),
				),
			),
		),
		ErrCompleted,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](0, 2),
			rx.SkipUntil[string](
				rx.Pipe(
					rx.Throw[int](ErrTest),
					DelaySubscription[int](3),
				),
			),
		),
		ErrTest,
	)
}
