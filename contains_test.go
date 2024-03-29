package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestContains(t *testing.T) {
	t.Parallel()

	greaterThanFour := func(v int) bool {
		return v > 4
	}

	NewTestSuite[bool](t).Case(
		rx.Pipe1(
			rx.Range(1, 10),
			rx.Contains(greaterThanFour),
		),
		true, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Range(1, 5),
			rx.Contains(greaterThanFour),
		),
		false, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 10),
				rx.Throw[int](ErrTest),
			),
			rx.Contains(greaterThanFour),
		),
		true, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 5),
				rx.Throw[int](ErrTest),
			),
			rx.Contains(greaterThanFour),
		),
		ErrTest,
	).Case(
		rx.Pipe2(
			rx.Range(1, 10),
			rx.Contains(greaterThanFour),
			rx.DoOnNext(func(bool) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe2(
			rx.Range(1, 5),
			rx.Contains(greaterThanFour),
			rx.DoOnNext(func(bool) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}

func TestContainsElement(t *testing.T) {
	t.Parallel()

	NewTestSuite[bool](t).Case(
		rx.Pipe1(
			rx.Range(1, 10),
			rx.ContainsElement(5),
		),
		true, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Range(1, 5),
			rx.ContainsElement(5),
		),
		false, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 10),
				rx.Throw[int](ErrTest),
			),
			rx.ContainsElement(5),
		),
		true, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 5),
				rx.Throw[int](ErrTest),
			),
			rx.ContainsElement(5),
		),
		ErrTest,
	).Case(
		rx.Pipe2(
			rx.Range(1, 10),
			rx.ContainsElement(5),
			rx.DoOnNext(func(bool) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe2(
			rx.Range(1, 5),
			rx.ContainsElement(5),
			rx.DoOnNext(func(bool) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}
