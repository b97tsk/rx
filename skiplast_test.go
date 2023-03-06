package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestSkipLast(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe1(
			rx.Range(1, 7),
			rx.SkipLast[int](0),
		),
		1, 2, 3, 4, 5, 6, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Range(1, 7),
			rx.SkipLast[int](3),
		),
		1, 2, 3, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Just(1),
			rx.SkipLast[int](3),
		),
		ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Empty[int](),
			rx.SkipLast[int](3),
		),
		ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 7),
				rx.Throw[int](ErrTest),
			),
			rx.SkipLast[int](0),
		),
		1, 2, 3, 4, 5, 6, ErrTest,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 7),
				rx.Throw[int](ErrTest),
			),
			rx.SkipLast[int](3),
		),
		1, 2, 3, ErrTest,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just(1),
				rx.Throw[int](ErrTest),
			),
			rx.SkipLast[int](3),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Empty[int](),
				rx.Throw[int](ErrTest),
			),
			rx.SkipLast[int](3),
		),
		ErrTest,
	)
}
