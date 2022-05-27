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
		rx.Pipe(
			rx.Range(1, 10),
			rx.Contains(greaterThanFour),
		),
		true, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Range(1, 5),
			rx.Contains(greaterThanFour),
		),
		false, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Concat(
				rx.Range(1, 10),
				rx.Throw[int](ErrTest),
			),
			rx.Contains(greaterThanFour),
		),
		true, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Concat(
				rx.Range(1, 5),
				rx.Throw[int](ErrTest),
			),
			rx.Contains(greaterThanFour),
		),
		ErrTest,
	)
}

func TestContainsElement(t *testing.T) {
	t.Parallel()

	NewTestSuite[bool](t).Case(
		rx.Pipe(
			rx.Range(1, 10),
			rx.ContainsElement(5),
		),
		true, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Range(1, 5),
			rx.ContainsElement(5),
		),
		false, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Concat(
				rx.Range(1, 10),
				rx.Throw[int](ErrTest),
			),
			rx.ContainsElement(5),
		),
		true, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Concat(
				rx.Range(1, 5),
				rx.Throw[int](ErrTest),
			),
			rx.ContainsElement(5),
		),
		ErrTest,
	)
}
