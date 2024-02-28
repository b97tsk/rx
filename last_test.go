package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestLast(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe1(
			rx.Empty[int](),
			rx.Last[int](),
		),
		rx.ErrEmpty,
	).Case(
		rx.Pipe1(
			rx.Throw[int](ErrTest),
			rx.Last[int](),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.Just(1),
			rx.Last[int](),
		),
		1, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Just(1, 2),
			rx.Last[int](),
		),
		2, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just(1),
				rx.Throw[int](ErrTest),
			),
			rx.Last[int](),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just(1, 2),
				rx.Throw[int](ErrTest),
			),
			rx.Last[int](),
		),
		ErrTest,
	).Case(
		rx.Pipe2(
			rx.Just(1, 2),
			rx.Last[int](),
			rx.OnNext(func(int) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}

func TestLastOrElse(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe1(
			rx.Empty[int](),
			rx.LastOrElse(404),
		),
		404, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Throw[int](ErrTest),
			rx.LastOrElse(404),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.Just(1),
			rx.LastOrElse(404),
		),
		1, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Just(1, 2),
			rx.LastOrElse(404),
		),
		2, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just(1),
				rx.Throw[int](ErrTest),
			),
			rx.LastOrElse(404),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just(1, 2),
				rx.Throw[int](ErrTest),
			),
			rx.LastOrElse(404),
		),
		ErrTest,
	).Case(
		rx.Pipe2(
			rx.Empty[int](),
			rx.LastOrElse(404),
			rx.OnNext(func(int) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}
