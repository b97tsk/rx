package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestTimeout(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Just("A", "B", "C").Pipe(
			AddLatencyToValues(1, 1),
			operators.Timeout(Step(2)),
		),
		"A", "B", "C", Completed,
	).Case(
		rx.Just("A", "B", "C").Pipe(
			AddLatencyToValues(1, 3),
			operators.Timeout(Step(2)),
		),
		"A", rx.ErrTimeout,
	).TestAll()
}
