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
		rx.Pipe(
			rx.Range(1, 10),
			rx.Every(lessThanFive),
		),
		false, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Range(1, 5),
			rx.Every(lessThanFive),
		),
		true, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Concat(
				rx.Range(1, 10),
				rx.Throw[int](ErrTest),
			),
			rx.Every(lessThanFive),
		),
		false, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Concat(
				rx.Range(1, 5),
				rx.Throw[int](ErrTest),
			),
			rx.Every(lessThanFive),
		),
		ErrTest,
	)
}