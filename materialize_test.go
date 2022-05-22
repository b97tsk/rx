package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestMaterialize(t *testing.T) {
	t.Parallel()

	count := func(n int, _ rx.Notification[string]) int {
		return n + 1
	}

	NewTestSuite[int](t).Case(
		rx.Pipe2(
			rx.Empty[string](),
			rx.Materialize[string](),
			rx.Reduce(0, count),
		),
		1, ErrCompleted,
	).Case(
		rx.Pipe2(
			rx.Throw[string](ErrTest),
			rx.Materialize[string](),
			rx.Reduce(0, count),
		),
		1, ErrCompleted,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			rx.Materialize[string](),
			rx.Reduce(0, count),
		),
		4, ErrCompleted,
	).Case(
		rx.Pipe2(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Throw[string](ErrTest),
			),
			rx.Materialize[string](),
			rx.Reduce(0, count),
		),
		4, ErrCompleted,
	)
}
