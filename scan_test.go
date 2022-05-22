package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestScan(t *testing.T) {
	t.Parallel()

	max := func(v1, v2 int) int {
		if v1 > v2 {
			return v1
		}

		return v2
	}

	sum := func(v1, v2 int) int {
		return v1 + v2
	}

	NewTestSuite[int](t).Case(
		rx.Pipe(
			rx.Range(1, 7),
			rx.Scan(0, max),
		),
		1, 2, 3, 4, 5, 6, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Just(42),
			rx.Scan(0, max),
		),
		42, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Empty[int](),
			rx.Scan(0, max),
		),
		ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Range(1, 7),
			rx.Scan(0, sum),
		),
		1, 3, 6, 10, 15, 21, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Just(42),
			rx.Scan(0, sum),
		),
		42, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Empty[int](),
			rx.Scan(0, sum),
		),
		ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Throw[int](ErrTest),
			rx.Scan(0, sum),
		),
		ErrTest,
	)
}
