package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestDefaultIfEmpty(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe1(
			rx.Empty[int](),
			rx.DefaultIfEmpty(42),
		),
		42, ErrCompleted,
	).Case(
		rx.Pipe1(
			rx.Range(1, 4),
			rx.DefaultIfEmpty(42),
		),
		1, 2, 3, ErrCompleted,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 4),
				rx.Throw[int](ErrTest),
			),
			rx.DefaultIfEmpty(42),
		),
		1, 2, 3, ErrTest,
	)
}

func TestThrowIfEmpty(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe1(
			rx.Empty[int](),
			rx.ThrowIfEmpty[int](),
		),
		rx.ErrEmpty,
	).Case(
		rx.Pipe1(
			rx.Range(1, 4),
			rx.ThrowIfEmpty[int](),
		),
		1, 2, 3, ErrCompleted,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 4),
				rx.Throw[int](ErrTest),
			),
			rx.ThrowIfEmpty[int](),
		),
		1, 2, 3, ErrTest,
	)
}

func TestSwitchIfEmpty(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe1(
			rx.Empty[int](),
			rx.SwitchIfEmpty(rx.Just(42)),
		),
		42, ErrCompleted,
	).Case(
		rx.Pipe1(
			rx.Range(1, 4),
			rx.SwitchIfEmpty(rx.Just(42)),
		),
		1, 2, 3, ErrCompleted,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 4),
				rx.Throw[int](ErrTest),
			),
			rx.SwitchIfEmpty(rx.Just(42)),
		),
		1, 2, 3, ErrTest,
	)
}
