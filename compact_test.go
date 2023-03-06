package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestCompact(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe1(
			rx.Just(1, 2, 2, 1, 3, 3, 1),
			rx.CompactComparable[int](),
		),
		1, 2, 1, 3, 1, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Just(1, 2, 2, 1, 3, 3, 1),
			rx.Compact(
				func(v1, v2 int) bool { return v1&1 == v2&1 },
			),
		),
		1, 2, 1, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Just(1, 2, 2, 1, 3, 3, 1),
			rx.CompactComparableKey(
				func(v int) int { return v & 1 },
			),
		),
		1, 2, 1, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Just(1, 2, 2, 1, 3, 3, 1),
			rx.CompactKey(
				func(v int) int { return v & 1 },
				func(v1, v2 int) bool { return v1 == v2 },
			),
		),
		1, 2, 1, ErrComplete,
	)
}
