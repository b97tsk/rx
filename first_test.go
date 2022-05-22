package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestFirst(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe(
			rx.Empty[int](),
			rx.First[int](),
		),
		rx.ErrEmpty,
	).Case(
		rx.Pipe(
			rx.Throw[int](ErrTest),
			rx.First[int](),
		),
		ErrTest,
	).Case(
		rx.Pipe(
			rx.Just(1),
			rx.First[int](),
		),
		1, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Just(1, 2),
			rx.First[int](),
		),
		1, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Concat(
				rx.Just(1),
				rx.Throw[int](ErrTest),
			),
			rx.First[int](),
		),
		1, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Concat(
				rx.Just(1, 2),
				rx.Throw[int](ErrTest),
			),
			rx.First[int](),
		),
		1, ErrCompleted,
	)
}

func TestFirstOrDefault(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe(
			rx.Empty[int](),
			rx.FirstOrDefault(404),
		),
		404, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Throw[int](ErrTest),
			rx.FirstOrDefault(404),
		),
		ErrTest,
	).Case(
		rx.Pipe(
			rx.Just(1),
			rx.FirstOrDefault(404),
		),
		1, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Just(1, 2),
			rx.FirstOrDefault(404),
		),
		1, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Concat(
				rx.Just(1),
				rx.Throw[int](ErrTest),
			),
			rx.FirstOrDefault(404),
		),
		1, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Concat(
				rx.Just(1, 2),
				rx.Throw[int](ErrTest),
			),
			rx.FirstOrDefault(404),
		),
		1, ErrCompleted,
	)
}
