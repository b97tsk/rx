package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestTake(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe(
			rx.Range(1, 10),
			rx.Take[int](0),
		),
		ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Range(1, 10),
			rx.Take[int](3),
		),
		1, 2, 3, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Just(1),
			rx.Take[int](3),
		),
		1, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Empty[int](),
			rx.Take[int](3),
		),
		ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Concat(
				rx.Range(1, 10),
				rx.Throw[int](ErrTest),
			),
			rx.Take[int](0),
		),
		ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Concat(
				rx.Range(1, 10),
				rx.Throw[int](ErrTest),
			),
			rx.Take[int](3),
		),
		1, 2, 3, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Concat(
				rx.Just(1),
				rx.Throw[int](ErrTest),
			),
			rx.Take[int](3),
		),
		1, ErrTest,
	).Case(
		rx.Pipe(
			rx.Concat(
				rx.Empty[int](),
				rx.Throw[int](ErrTest),
			),
			rx.Take[int](3),
		),
		ErrTest,
	)
}
