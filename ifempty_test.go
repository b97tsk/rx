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
			rx.DefaultIfEmpty(1, 2, 3),
		),
		1, 2, 3, ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Empty[int](),
			rx.DefaultIfEmpty(1, 2, 3),
			rx.Take[int](1),
		),
		1, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Range(1, 4),
			rx.DefaultIfEmpty(4, 5, 6),
		),
		1, 2, 3, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 4),
				rx.Throw[int](ErrTest),
			),
			rx.DefaultIfEmpty(4, 5, 6),
		),
		1, 2, 3, ErrTest,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 4),
				rx.Oops[int](ErrTest),
			),
			rx.DefaultIfEmpty(4, 5, 6),
		),
		1, 2, 3, rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe2(
			rx.Empty[int](),
			rx.DefaultIfEmpty(1, 2, 3),
			rx.DoOnNext(func(int) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
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
		1, 2, 3, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 4),
				rx.Throw[int](ErrTest),
			),
			rx.ThrowIfEmpty[int](),
		),
		1, 2, 3, ErrTest,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 4),
				rx.Oops[int](ErrTest),
			),
			rx.ThrowIfEmpty[int](),
		),
		1, 2, 3, rx.ErrOops, ErrTest,
	)
}

func TestSwitchIfEmpty(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe1(
			rx.Empty[int](),
			rx.SwitchIfEmpty(rx.Just(1, 2, 3)),
		),
		1, 2, 3, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Range(1, 4),
			rx.SwitchIfEmpty(rx.Just(4, 5, 6)),
		),
		1, 2, 3, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 4),
				rx.Throw[int](ErrTest),
			),
			rx.SwitchIfEmpty(rx.Just(4, 5, 6)),
		),
		1, 2, 3, ErrTest,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 4),
				rx.Oops[int](ErrTest),
			),
			rx.SwitchIfEmpty(rx.Just(4, 5, 6)),
		),
		1, 2, 3, rx.ErrOops, ErrTest,
	)
}
