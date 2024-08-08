package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestTake(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe1(
			rx.Range(1, 10),
			rx.Take[int](0),
		),
		ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Range(1, 10),
			rx.Take[int](3),
		),
		1, 2, 3, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Just(1),
			rx.Take[int](3),
		),
		1, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Empty[int](),
			rx.Take[int](3),
		),
		ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 10),
				rx.Throw[int](ErrTest),
			),
			rx.Take[int](0),
		),
		ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 10),
				rx.Throw[int](ErrTest),
			),
			rx.Take[int](3),
		),
		1, 2, 3, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just(1),
				rx.Throw[int](ErrTest),
			),
			rx.Take[int](3),
		),
		1, ErrTest,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Empty[int](),
				rx.Throw[int](ErrTest),
			),
			rx.Take[int](3),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 10),
				rx.Oops[int](ErrTest),
			),
			rx.Take[int](0),
		),
		ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 10),
				rx.Oops[int](ErrTest),
			),
			rx.Take[int](3),
		),
		1, 2, 3, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just(1),
				rx.Oops[int](ErrTest),
			),
			rx.Take[int](3),
		),
		1, rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Empty[int](),
				rx.Oops[int](ErrTest),
			),
			rx.Take[int](3),
		),
		rx.ErrOops, ErrTest,
	)
}
