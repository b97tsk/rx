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

func TestFirstOrElse(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe1(
			rx.Empty[int](),
			rx.FirstOrElse(404),
		),
		404, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Throw[int](ErrTest),
			rx.FirstOrElse(404),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.Just(1),
			rx.FirstOrElse(404),
		),
		1, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Just(1, 2),
			rx.FirstOrElse(404),
		),
		1, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just(1),
				rx.Throw[int](ErrTest),
			),
			rx.FirstOrElse(404),
		),
		1, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just(1, 2),
				rx.Throw[int](ErrTest),
			),
			rx.FirstOrElse(404),
		),
		1, ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Empty[int](),
			rx.FirstOrElse(404),
			rx.DoOnNext(func(int) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}
