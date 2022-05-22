package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestLast(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe(
			rx.Empty[int](),
			rx.Last[int](),
		),
		rx.ErrEmpty,
	).Case(
		rx.Pipe(
			rx.Throw[int](ErrTest),
			rx.Last[int](),
		),
		ErrTest,
	).Case(
		rx.Pipe(
			rx.Just(1),
			rx.Last[int](),
		),
		1, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Just(1, 2),
			rx.Last[int](),
		),
		2, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Concat(
				rx.Just(1),
				rx.Throw[int](ErrTest),
			),
			rx.Last[int](),
		),
		ErrTest,
	).Case(
		rx.Pipe(
			rx.Concat(
				rx.Just(1, 2),
				rx.Throw[int](ErrTest),
			),
			rx.Last[int](),
		),
		ErrTest,
	)
}

func TestLastOrDefault(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe(
			rx.Empty[int](),
			rx.LastOrDefault(404),
		),
		404, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Throw[int](ErrTest),
			rx.LastOrDefault(404),
		),
		ErrTest,
	).Case(
		rx.Pipe(
			rx.Just(1),
			rx.LastOrDefault(404),
		),
		1, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Just(1, 2),
			rx.LastOrDefault(404),
		),
		2, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Concat(
				rx.Just(1),
				rx.Throw[int](ErrTest),
			),
			rx.LastOrDefault(404),
		),
		ErrTest,
	).Case(
		rx.Pipe(
			rx.Concat(
				rx.Just(1, 2),
				rx.Throw[int](ErrTest),
			),
			rx.LastOrDefault(404),
		),
		ErrTest,
	)
}
