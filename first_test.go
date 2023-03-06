package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestFirst(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe1(
			rx.Empty[int](),
			rx.First[int](),
		),
		rx.ErrEmpty,
	).Case(
		rx.Pipe1(
			rx.Throw[int](ErrTest),
			rx.First[int](),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.Just(1),
			rx.First[int](),
		),
		1, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Just(1, 2),
			rx.First[int](),
		),
		1, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just(1),
				rx.Throw[int](ErrTest),
			),
			rx.First[int](),
		),
		1, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just(1, 2),
				rx.Throw[int](ErrTest),
			),
			rx.First[int](),
		),
		1, ErrComplete,
	)
}

func TestFirstOrDefault(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe1(
			rx.Empty[int](),
			rx.FirstOrDefault(404),
		),
		404, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Throw[int](ErrTest),
			rx.FirstOrDefault(404),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.Just(1),
			rx.FirstOrDefault(404),
		),
		1, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Just(1, 2),
			rx.FirstOrDefault(404),
		),
		1, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just(1),
				rx.Throw[int](ErrTest),
			),
			rx.FirstOrDefault(404),
		),
		1, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just(1, 2),
				rx.Throw[int](ErrTest),
			),
			rx.FirstOrDefault(404),
		),
		1, ErrComplete,
	)
}
