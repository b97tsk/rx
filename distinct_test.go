package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestDistinct(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe1(
			rx.Just(1, 2, 2, 1, 3, 3, 1),
			rx.DistinctComparable[int](),
		),
		1, 2, 3, ErrCompleted,
	).Case(
		rx.Pipe1(
			rx.Just(1, 2, 2, 1, 3, 3, 1),
			rx.Distinct(
				func(v int) int { return v & 1 },
			),
		),
		1, 2, ErrCompleted,
	)
}
