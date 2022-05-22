package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestSkip(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe(
			rx.Range(1, 7),
			rx.Skip[int](0),
		),
		1, 2, 3, 4, 5, 6, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Range(1, 7),
			rx.Skip[int](3),
		),
		4, 5, 6, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Just(1),
			rx.Skip[int](3),
		),
		ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Empty[int](),
			rx.Skip[int](3),
		),
		ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Concat(
				rx.Range(1, 7),
				rx.Throw[int](ErrTest),
			),
			rx.Skip[int](0),
		),
		1, 2, 3, 4, 5, 6, ErrTest,
	).Case(
		rx.Pipe(
			rx.Concat(
				rx.Range(1, 7),
				rx.Throw[int](ErrTest),
			),
			rx.Skip[int](3),
		),
		4, 5, 6, ErrTest,
	).Case(
		rx.Pipe(
			rx.Concat(
				rx.Just(1),
				rx.Throw[int](ErrTest),
			),
			rx.Skip[int](3),
		),
		ErrTest,
	).Case(
		rx.Pipe(
			rx.Concat(
				rx.Empty[int](),
				rx.Throw[int](ErrTest),
			),
			rx.Skip[int](3),
		),
		ErrTest,
	)
}
