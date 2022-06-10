package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestPairwise(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe2(
			rx.Empty[string](),
			rx.PairWise[string](),
			ToString[rx.Pair[string, string]](),
		),
		ErrCompleted,
	).Case(
		rx.Pipe2(
			rx.Just("A"),
			rx.PairWise[string](),
			ToString[rx.Pair[string, string]](),
		),
		ErrCompleted,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B"),
			rx.PairWise[string](),
			ToString[rx.Pair[string, string]](),
		),
		"{A B}", ErrCompleted,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			rx.PairWise[string](),
			ToString[rx.Pair[string, string]](),
		),
		"{A B}", "{B C}", ErrCompleted,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C", "D"),
			rx.PairWise[string](),
			ToString[rx.Pair[string, string]](),
		),
		"{A B}", "{B C}", "{C D}", ErrCompleted,
	).Case(
		rx.Pipe2(
			rx.Concat(
				rx.Just("A", "B", "C", "D"),
				rx.Throw[string](ErrTest),
			),
			rx.PairWise[string](),
			ToString[rx.Pair[string, string]](),
		),
		"{A B}", "{B C}", "{C D}", ErrTest,
	)
}