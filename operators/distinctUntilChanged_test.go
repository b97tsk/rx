package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestDistinctUntilChanged(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Just("A", "B", "B", "A", "C", "C", "A").Pipe(
			operators.DistinctUntilChanged(),
		),
		"A", "B", "A", "C", "A", Completed,
	).TestAll()
}
