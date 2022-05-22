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
		1, 2, 3, 4, 5, 6, 7, 8, 9, ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Range(1, 10),
			rx.Take[int](5),
		),
		1, 2, 3, 4, 5, ErrCompleted,
	)
}

func TestIota(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe(
			rx.Iota(1),
			rx.Take[int](5),
		),
		1, 2, 3, 4, 5, ErrCompleted,
	)
}
