package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestMerge(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Merge[string](),
		ErrCompleted,
	).Case(
		rx.Merge(
			rx.Pipe(rx.Just("A", "B"), AddLatencyToValues[string](3, 5)),
			rx.Pipe(rx.Just("C", "D"), AddLatencyToValues[string](2, 4)),
			rx.Pipe(rx.Just("E", "F"), AddLatencyToValues[string](1, 3)),
		),
		"E", "C", "A", "F", "D", "B", ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Pipe(rx.Just("A", "B"), AddLatencyToValues[string](3, 5)),
			rx.MergeWith(
				rx.Pipe(rx.Just("C", "D"), AddLatencyToValues[string](2, 4)),
				rx.Pipe(rx.Just("E", "F"), AddLatencyToValues[string](1, 3)),
			),
		),
		"E", "C", "A", "F", "D", "B", ErrCompleted,
	)
}
