package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestMergeAll(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Just(
			rx.Just("A", "B").Pipe(AddLatencyToValues(3, 5)),
			rx.Just("C", "D").Pipe(AddLatencyToValues(2, 4)),
			rx.Just("E", "F").Pipe(AddLatencyToValues(1, 3)),
		).Pipe(
			operators.MergeAll(-1),
		),
		"E", "C", "A", "F", "D", "B", Completed,
	).TestAll()
}
