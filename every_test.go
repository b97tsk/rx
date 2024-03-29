package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestEvery(t *testing.T) {
	t.Parallel()

	lessThanFive := func(v int) bool {
		return v < 5
	}

	NewTestSuite[bool](t).Case(
		rx.Pipe1(
			rx.Range(1, 10),
			rx.Every(lessThanFive),
		),
		false, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Range(1, 5),
			rx.Every(lessThanFive),
		),
		true, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 10),
				rx.Throw[int](ErrTest),
			),
			rx.Every(lessThanFive),
		),
		false, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 5),
				rx.Throw[int](ErrTest),
			),
			rx.Every(lessThanFive),
		),
		ErrTest,
	).Case(
		rx.Pipe2(
			rx.Range(1, 10),
			rx.Every(lessThanFive),
			rx.DoOnNext(func(bool) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe2(
			rx.Range(1, 5),
			rx.Every(lessThanFive),
			rx.DoOnNext(func(bool) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}
