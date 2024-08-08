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
		rx.Pipe1(
			rx.Range(1, 7),
			rx.Scan(0, max),
		),
		1, 2, 3, 4, 5, 6, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Just(42),
			rx.Scan(0, max),
		),
		42, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Empty[int](),
			rx.Scan(0, max),
		),
		ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Range(1, 7),
			rx.Scan(0, sum),
		),
		1, 3, 6, 10, 15, 21, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Just(42),
			rx.Scan(0, sum),
		),
		42, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Empty[int](),
			rx.Scan(0, sum),
		),
		ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Throw[int](ErrTest),
			rx.Scan(0, sum),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.Oops[int](ErrTest),
			rx.Scan(0, sum),
		),
		rx.ErrOops, ErrTest,
	)
}
