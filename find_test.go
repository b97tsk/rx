package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestFind(t *testing.T) {
	t.Parallel()

	equalFive := func(v int) bool {
		return v == 5
	}

	NewTestSuite[int](t).Case(
		rx.Pipe1(
			rx.Range(1, 10),
			rx.Find(equalFive),
		),
		5, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Range(1, 5),
			rx.Find(equalFive),
		),
		ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 10),
				rx.Throw[int](ErrTest),
			),
			rx.Find(equalFive),
		),
		5, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 5),
				rx.Throw[int](ErrTest),
			),
			rx.Find(equalFive),
		),
		ErrTest,
	)
}
