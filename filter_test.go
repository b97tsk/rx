package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestFilter(t *testing.T) {
	t.Parallel()

	lessThanFive := func(v int) bool {
		return v < 5
	}

	NewTestSuite[int](t).Case(
		rx.Pipe1(
			rx.Range(1, 10),
			rx.Filter(lessThanFive),
		),
		1, 2, 3, 4, ErrCompleted,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 10),
				rx.Throw[int](ErrTest),
			),
			rx.Filter(lessThanFive),
		),
		1, 2, 3, 4, ErrTest,
	)
}

func TestFilterOut(t *testing.T) {
	t.Parallel()

	lessThanFive := func(v int) bool {
		return v < 5
	}

	NewTestSuite[int](t).Case(
		rx.Pipe1(
			rx.Range(1, 10),
			rx.FilterOut(lessThanFive),
		),
		5, 6, 7, 8, 9, ErrCompleted,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 10),
				rx.Throw[int](ErrTest),
			),
			rx.FilterOut(lessThanFive),
		),
		5, 6, 7, 8, 9, ErrTest,
	)
}

func TestFilterMap(t *testing.T) {
	t.Parallel()

	lessThanFive := func(v int) (int, bool) {
		return v * 2, v < 5
	}

	NewTestSuite[int](t).Case(
		rx.Pipe1(
			rx.Range(1, 10),
			rx.FilterMap(lessThanFive),
		),
		2, 4, 6, 8, ErrCompleted,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 10),
				rx.Throw[int](ErrTest),
			),
			rx.FilterMap(lessThanFive),
		),
		2, 4, 6, 8, ErrTest,
	)
}
