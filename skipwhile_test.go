package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestSkipWhile(t *testing.T) {
	t.Parallel()

	lessThanFive := func(v int) bool {
		return v < 5
	}

	NewTestSuite[int](t).Case(
		rx.Pipe1(
			rx.Just(1, 2, 3, 4, 5, 4, 3, 2, 1),
			rx.SkipWhile(lessThanFive),
		),
		5, 4, 3, 2, 1, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 10),
				rx.Throw[int](ErrTest),
			),
			rx.SkipWhile(lessThanFive),
		),
		5, 6, 7, 8, 9, ErrTest,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 5),
				rx.Throw[int](ErrTest),
			),
			rx.SkipWhile(lessThanFive),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 5),
				rx.Oops[int](ErrTest),
			),
			rx.SkipWhile(lessThanFive),
		),
		rx.ErrOops, ErrTest,
	)
}
