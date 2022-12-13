package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestReduce(t *testing.T) {
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
			rx.Reduce(0, max),
		),
		6, ErrCompleted,
	).Case(
		rx.Pipe1(
			rx.Just(42),
			rx.Reduce(0, max),
		),
		42, ErrCompleted,
	).Case(
		rx.Pipe1(
			rx.Empty[int](),
			rx.Reduce(0, max),
		),
		0, ErrCompleted,
	).Case(
		rx.Pipe1(
			rx.Range(1, 7),
			rx.Reduce(0, sum),
		),
		21, ErrCompleted,
	).Case(
		rx.Pipe1(
			rx.Just(42),
			rx.Reduce(0, sum),
		),
		42, ErrCompleted,
	).Case(
		rx.Pipe1(
			rx.Empty[int](),
			rx.Reduce(0, sum),
		),
		0, ErrCompleted,
	).Case(
		rx.Pipe1(
			rx.Throw[int](ErrTest),
			rx.Reduce(0, sum),
		),
		ErrTest,
	)
}
