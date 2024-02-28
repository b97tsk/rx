package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestRange(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Range(1, 10),
		1, 2, 3, 4, 5, 6, 7, 8, 9, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Range(1, 10),
			rx.Take[int](5),
		),
		1, 2, 3, 4, 5, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Range(1, 10),
			rx.OnNext(func(int) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}

func TestIota(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe1(
			rx.Iota(1),
			rx.Take[int](5),
		),
		1, 2, 3, 4, 5, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Iota(1),
			rx.OnNext(func(int) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}
