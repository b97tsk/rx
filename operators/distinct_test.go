package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestDistinct(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Just("A", "B", "B", "A", "C", "C", "A").Pipe(
			operators.Distinct(),
		),
		"A", "B", "C", Completed,
	).TestAll()
}
